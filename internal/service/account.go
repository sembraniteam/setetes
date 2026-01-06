package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/sembraniteam/setetes/internal"
	"github.com/sembraniteam/setetes/internal/cryptox"
	"github.com/sembraniteam/setetes/internal/cryptox/argon2x"
	"github.com/sembraniteam/setetes/internal/cryptox/pasetox"
	"github.com/sembraniteam/setetes/internal/ent"
	"github.com/sembraniteam/setetes/internal/ent/account"
	"github.com/sembraniteam/setetes/internal/ent/otp"
	"github.com/sembraniteam/setetes/internal/httpx/request"
)

const (
	exp     = time.Minute * 30
	skew    = time.Second * -30
	charLen = 6
)

type (
	AccountQuery struct {
		client *ent.Client
		ctx    context.Context
	}

	Account interface {
		Authorize(body request.Authorize) (*pasetox.TokenPair, error)
		Register(body request.Account) error
		Activate(body request.Activation) error
		ResendOTP(body request.ResendOTP) error
		Self(id uuid.UUID) (*ent.Account, error)
	}
)

func NewAccount(i do.Injector) (Account, error) {
	return &AccountQuery{
		client: do.MustInvoke[*ent.Client](i),
		ctx:    context.Background(),
	}, nil
}

func (a *AccountQuery) Authorize(
	body request.Authorize,
) (*pasetox.TokenPair, error) {
	acc, err := a.client.Account.Query().
		Where(
			account.EmailEQ(body.Email),
			account.LockedEQ(false),
			account.ActivatedEQ(true),
		).
		WithPassword().
		Only(a.ctx)
	if err != nil {
		return nil, err
	}

	ac := argon2x.Default()
	_, err = ac.VerifyString(
		[]byte(body.Password),
		acc.Edges.Password.Hash,
	)
	if err != nil {
		return nil, err
	}

	tokenPair, err := generateToken(acc.ID, body.Platform)
	if err != nil {
		return nil, err
	}

	return tokenPair, nil
}

func (a *AccountQuery) Register(body request.Account) error {
	tx, err := a.client.Tx(a.ctx)
	if err != nil {
		return err
	}

	acc, err := tx.Account.Create().
		SetNationalIDHash(cryptox.Sha256(body.NationalID)).
		SetNationalIDMasked(cryptox.MaskNumber(body.NationalID)).
		SetFullName(body.FullName).
		SetGender(body.GetGender()).
		SetEmail(body.Email).
		SetCountryIsoCode(body.CountryISOCode).
		SetDialCode(body.DialCode).
		SetPhoneNumber(body.PhoneNumber).
		Save(a.ctx)
	if err != nil {
		return rollback(tx, err)
	}

	code, err := genOTP()
	if err != nil {
		return rollback(tx, err)
	}

	// TODO: send OTP code to email.
	print("The OTP Code is ", code)

	_, err = tx.OTP.Create().
		SetCodeHash(cryptox.Sha256(code)).
		SetType(otp.TypeActivation).
		SetAccount(acc).
		SetExpiredAt(time.Now().Add(exp).UnixMilli()).Save(a.ctx)
	if err != nil {
		return rollback(tx, err)
	}

	return tx.Commit()
}

func (a *AccountQuery) Activate(body request.Activation) error {
	tx, err := a.client.Tx(a.ctx)
	if err != nil {
		return err
	}

	otps, err := tx.OTP.Query().
		Where(
			otp.TypeEQ(otp.TypeActivation),
			otp.ExpiredAtGTE(time.Now().UnixMilli()),
			otp.HasAccountWith(account.Not(account.HasPassword())),
		).
		WithAccount().
		All(a.ctx)
	if err != nil {
		return rollback(tx, err)
	}

	var validOtp *ent.OTP
	for _, o := range otps {
		if cryptox.VerifySha256(body.Code, o.CodeHash) {
			validOtp = o
			break
		}
	}

	if validOtp == nil {
		return rollback(tx, errors.New("invalid or expired OTP"))
	}

	acc, err := validOtp.QueryAccount().Only(a.ctx)
	if err != nil {
		return rollback(tx, err)
	}

	if err = tx.OTP.DeleteOne(validOtp).Exec(a.ctx); err != nil {
		return rollback(tx, err)
	}

	config := argon2x.Default()
	pwd, err := config.HashString([]byte(body.Password))
	if err != nil {
		return rollback(tx, err)
	}

	_, err = tx.Password.Create().
		SetHash(pwd).
		SetAccount(acc).
		Save(a.ctx)
	if err != nil {
		return rollback(tx, err)
	}

	if err = tx.Account.UpdateOne(acc).
		SetActivated(true).
		Exec(a.ctx); err != nil {
		return rollback(tx, err)
	}

	return tx.Commit()
}

func (a *AccountQuery) ResendOTP(body request.ResendOTP) error {
	tx, err := a.client.Tx(a.ctx)
	if err != nil {
		return err
	}

	acc, err := tx.Account.Query().
		Where(
			account.EmailEQ(body.Email),
			account.LockedEQ(false),
		).
		Only(a.ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return rollback(
				tx,
				errors.New("account not found"),
			)
		}

		return rollback(tx, err)
	}

	code, err := genOTP()
	if err != nil {
		return rollback(tx, err)
	}

	// TODO: send OTP code to email.
	print("The OTP Code is ", code)

	_, err = tx.OTP.Create().
		SetCodeHash(cryptox.Sha256(code)).
		SetType(body.GetType()).
		SetAccount(acc).
		SetExpiredAt(time.Now().Add(exp).UnixMilli()).Save(a.ctx)
	if err != nil {
		return rollback(tx, err)
	}

	return tx.Commit()
}

func (a *AccountQuery) Self(id uuid.UUID) (*ent.Account, error) {
	return a.client.Account.Get(a.ctx, id)
}

func genOTP() (string, error) {
	return cryptox.RandChars(charLen)
}

func rollback(tx *ent.Tx, err error) error {
	if rerr := tx.Rollback(); rerr != nil {
		err = fmt.Errorf("%w: %v", err, rerr)
	}

	return err
}

func generateToken(
	subject uuid.UUID,
	platform string,
) (*pasetox.TokenPair, error) {
	ed := internal.Get().ED25519
	privateKey, err := cryptox.LoadPrivateKey(ed.PrivateKeyPath)
	if err != nil {
		return nil, err
	}

	publicKey, err := cryptox.LoadPublicKey(ed.PublicKeyPath)
	if err != nil {
		return nil, err
	}
	kp, err := cryptox.NewKeypair(privateKey, publicKey)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	token := pasetox.New(kp, pasetox.Claims{
		Platform:        platform,
		Subject:         subject.String(),
		TokenIdentifier: uuid.NewString(),
		Expiration:      now.Add(exp),
		IssuedAt:        now.Add(skew),
		NotBefore:       now.Add(skew),
	})

	accessToken, err := token.Signed()
	if err != nil {
		return nil, err
	}

	return &pasetox.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: cryptox.RandToken(),
		ExpiresIn:    now.Add(exp).UnixMilli(),
	}, nil
}

package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/samber/do/v2"
	"github.com/sembraniteam/setetes/internal/cryptox"
	"github.com/sembraniteam/setetes/internal/cryptox/argon2x"
	"github.com/sembraniteam/setetes/internal/ent"
	"github.com/sembraniteam/setetes/internal/ent/account"
	"github.com/sembraniteam/setetes/internal/ent/otp"
	"github.com/sembraniteam/setetes/internal/httpx/request"
)

const (
	exp     = time.Minute * 30
	charLen = 6
)

type (
	AccountQuery struct {
		client *ent.Client
		ctx    context.Context
	}

	Account interface {
		Register(body request.Account) error
		Activate(body request.Activation) error
		ResendOTP(body request.ResendOTP) error
	}
)

func NewAccount(i do.Injector) (Account, error) {
	return &AccountQuery{
		client: do.MustInvoke[*ent.Client](i),
		ctx:    context.Background(),
	}, nil
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

func genOTP() (string, error) {
	return cryptox.RandChars(charLen)
}

func rollback(tx *ent.Tx, err error) error {
	if rerr := tx.Rollback(); rerr != nil {
		err = fmt.Errorf("%w: %v", err, rerr)
	}

	return err
}

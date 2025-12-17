package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/megalodev/setetes/internal/cryptox"
	"github.com/megalodev/setetes/internal/cryptox/argon2x"
	"github.com/megalodev/setetes/internal/ent"
	"github.com/megalodev/setetes/internal/ent/otp"
	"github.com/megalodev/setetes/internal/httpx/request"
	"github.com/samber/do/v2"
)

const (
	expiredTime = time.Minute * 30
	charLen     = 6
)

type (
	AccountQuery struct {
		client *ent.Client
		ctx    context.Context
	}

	Account interface {
		Register(body request.Account) error
		Activate(body request.Activation) error
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

	account, err := tx.Account.Create().
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

	code, err := cryptox.RandChars(charLen)
	if err != nil {
		return rollback(tx, err)
	}

	// TODO: send OTP code to email.
	print("OTP Code is ", code)

	_, err = tx.OTP.Create().
		SetCodeHash(cryptox.Sha256(code)).
		SetType(otp.TypeRegister).
		SetAccount(account).
		SetExpiredAt(time.Now().Add(expiredTime).UnixMilli()).Save(a.ctx)
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
			otp.TypeEQ(otp.TypeRegister),
			otp.ExpiredAtGTE(time.Now().UnixMilli()),
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
		SetAccount(validOtp.Edges.Account).
		Save(a.ctx)
	if err != nil {
		return rollback(tx, err)
	}

	if err = tx.Account.UpdateOne(validOtp.Edges.Account).
		SetActivated(true).
		Exec(a.ctx); err != nil {
		return rollback(tx, err)
	}

	return tx.Commit()
}

func rollback(tx *ent.Tx, err error) error {
	if rerr := tx.Rollback(); rerr != nil {
		err = fmt.Errorf("%w: %v", err, rerr)
	}

	return err
}

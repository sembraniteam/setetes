package service

import (
	"context"
	"fmt"
	"time"

	"github.com/megalodev/setetes/internal/cryptox"
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

	_, err = tx.OTP.Create().
		SetCode(code).
		SetType(otp.TypeRegister).
		SetAccount(account).
		SetExpiredAt(time.Now().Add(expiredTime).UnixMilli()).Save(a.ctx)
	if err != nil {
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

package service

import (
	"context"

	"github.com/megalodev/setetes/internal/cryptox"
	"github.com/megalodev/setetes/internal/ent"
	"github.com/megalodev/setetes/internal/httpx/request"
	"github.com/samber/do/v2"
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
	_, err := a.client.Account.Create().
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
		return err
	}

	return nil
}

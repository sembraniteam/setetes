package service

import (
	"context"

	"github.com/megalodev/setetes/internal/ent"
	"github.com/samber/do/v2"
)

type (
	AccountQuery struct {
		client *ent.Client
	}

	Account interface {
		Register() error
	}
)

func NewAccount(i do.Injector) (Account, error) {
	return &AccountQuery{client: do.MustInvoke[*ent.Client](i)}, nil
}

func (a *AccountQuery) Register() error {
	_, err := a.client.Account.Create().
		SetFullName("").
		Save(context.Background())
	if err != nil {
		return err
	}

	return nil
}

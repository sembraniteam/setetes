package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/megalodev/setetes/internal/ent"
	"github.com/megalodev/setetes/internal/httpx"
	"github.com/megalodev/setetes/internal/httpx/request"
	"github.com/megalodev/setetes/internal/httpx/response"
	"github.com/megalodev/setetes/internal/service"
	"github.com/samber/do/v2"
)

type (
	Account struct {
		service service.Account
	}
)

func NewAccount(i do.Injector) (Account, error) {
	return Account{
		service: do.MustInvoke[service.Account](i),
	}, nil
}

func (a *Account) Register(ctx *gin.Context) {
	body, berr := response.ValidateJSON[request.Account](ctx)
	if berr != nil {
		response.Error(ctx, berr)
		return
	}

	if err := a.service.Register(*body); err != nil {
		if ent.IsConstraintError(err) {
			response.BadRequest(
				ctx,
				httpx.DuplicateKeyCode,
				response.NewMessage(
					response.Warning,
					"account is already registered",
				),
			)
			return
		}

		response.Error(ctx, err)
		return
	}

	response.Created(ctx, nil, nil)
}

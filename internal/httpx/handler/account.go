package handler

import (
	"log/slog"

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
		log     *slog.Logger
	}
)

func NewAccount(i do.Injector) (Account, error) {
	return Account{
		service: do.MustInvoke[service.Account](i),
		log:     slog.Default(),
	}, nil
}

func (a *Account) Register(ctx *gin.Context) {
	body, berr := response.ValidateJSON[request.Account](ctx)
	if berr != nil {
		a.log.Error("validate request failed", slog.Any("error", berr))
		response.Error(ctx, berr)
		return
	}

	if err := a.service.Register(*body); err != nil {
		a.log.Error("register account failed", slog.Any("error", err))
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

func (a *Account) Activate(ctx *gin.Context) {
	body, berr := response.ValidateJSON[request.Activation](ctx)
	if berr != nil {
		a.log.Error("validate request failed", slog.Any("error", berr))
		response.Error(ctx, berr)
		return
	}

	if err := a.service.Activate(*body); err != nil {
		a.log.Error("activate account failed", slog.Any("error", err))
		response.InvalidParameter(ctx, err.Error())
		return
	}

	response.Ok(ctx, response.MsgSuccess, nil)
}

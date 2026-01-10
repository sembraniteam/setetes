package handler

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"github.com/sembraniteam/setetes/internal/ent"
	"github.com/sembraniteam/setetes/internal/httpx"
	"github.com/sembraniteam/setetes/internal/httpx/request"
	"github.com/sembraniteam/setetes/internal/httpx/response"
	"github.com/sembraniteam/setetes/internal/httpx/response/responsetypes"
	"github.com/sembraniteam/setetes/internal/service"
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

func (a *Account) Authorize(ctx *gin.Context) {
	body, berr := response.ValidateJSON[request.Authorize](ctx)
	if berr != nil {
		a.log.Error("validate request failed", slog.Any("error", berr))
		response.Error(ctx, berr)
		return
	}

	tokenPair, err := a.service.Authorize(*body)
	if err != nil {
		a.log.Error("authorize failed", slog.Any("error", err))
		response.InvalidParameter(ctx, err.Error())
		return
	}

	response.Ok(ctx, response.MsgSuccess, tokenPair)
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

		response.InvalidParameter(ctx, err.Error())
		return
	}

	response.Created(ctx, response.MsgSuccess, nil)
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

func (a *Account) ResendOTP(ctx *gin.Context) {
	body, berr := response.ValidateJSON[request.ResendOTP](ctx)
	if berr != nil {
		a.log.Error("validate request failed", slog.Any("error", berr))
		response.Error(ctx, berr)
		return
	}

	if err := a.service.ResendOTP(*body); err != nil {
		a.log.Error("resend OTP failed", slog.Any("error", err))
		response.InvalidParameter(ctx, err.Error())
		return
	}

	response.Ok(ctx, response.MsgSuccess, nil)
}

func (a *Account) Self(ctx *gin.Context) {
	httpContext := httpx.NewContext(ctx)
	session := httpContext.GetUserSession()
	if session != nil {
		account, err := a.service.Self(session.ID)
		if err != nil {
			a.log.Error("get account failed", slog.Any("error", err))
			response.Error(ctx, err)
			return
		}

		acc := responsetypes.Account{Account: account}

		response.Ok(ctx, response.MsgSuccess, acc.ToResponse())
		return
	}
}

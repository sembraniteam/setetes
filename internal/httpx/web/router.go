package web

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"github.com/sembraniteam/setetes/internal/httpx/handler"
	"github.com/sembraniteam/setetes/internal/httpx/middleware"
	"github.com/sembraniteam/setetes/internal/httpx/response"
)

func Routes(e *gin.Engine, i do.Injector) {
	rateLimiter := middleware.NewTokenBucket(1, time.Minute*1)
	accountH := do.MustInvoke[handler.Account](i)

	e.GET("/ping", func(c *gin.Context) {
		response.Ok(c, response.MsgPong, nil)
	})

	accountG := e.Group("/account/v1")
	{
		accountG.POST("/authorization", accountH.Authorize)
		accountG.POST("/register", accountH.Register)
		accountG.POST("/activate", accountH.Activate)
		accountG.Use(middleware.RateLimitByIP(rateLimiter)).
			POST("/resend-otp", accountH.ResendOTP)
	}
}

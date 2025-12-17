package web

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
	"github.com/sembraniteam/setetes/internal/httpx/handler"
	"github.com/sembraniteam/setetes/internal/httpx/response"
)

func Routes(e *gin.Engine, i do.Injector) {
	accountH := do.MustInvoke[handler.Account](i)

	e.GET("/ping", func(c *gin.Context) {
		response.Ok(c, response.MsgPong, nil)
	})

	accountG := e.Group("/account/v1")
	{
		accountG.POST("/register", accountH.Register)
		accountG.POST("/activate", accountH.Activate)
	}
}

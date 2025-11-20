package web

import (
	"github.com/gin-gonic/gin"
	"github.com/megalodev/setetes/internal/httpx/response"
	"github.com/samber/do/v2"
)

func Routes(e *gin.Engine, i do.Injector) {
	v1 := e.Group("/v1")

	v1.GET("/ping", func(c *gin.Context) {
		response.Ok(c, response.MsgPong, nil)
	})
}

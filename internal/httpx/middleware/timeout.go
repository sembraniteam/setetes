package middleware

import (
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/megalodev/setetes/internal/httpx"
)

func Timeout() gin.HandlerFunc {
	return timeout.New(timeout.WithTimeout(httpx.Timeout))
}

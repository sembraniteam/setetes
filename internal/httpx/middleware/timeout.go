package middleware

import (
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
)

func Timeout() gin.HandlerFunc {
	return timeout.New(timeout.WithTimeout(time.Second * 30))
}

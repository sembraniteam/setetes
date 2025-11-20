package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/megalodev/setetes/internal/httpx"
	"github.com/megalodev/setetes/internal/httpx/response"
)

const headerXRequestID = "X-Request-ID"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader(headerXRequestID)
		if rid == "" {
			rid = uuid.NewString()
		}

		parsed, err := uuid.Parse(rid)
		if err != nil {
			response.BadRequest(c, httpx.InvalidValueCode, response.MsgInvalidRequestID.WithField(headerXRequestID))
			c.Abort()
			return
		}

		w := httpx.NewContext(c)
		w.SetRequestID(parsed)

		c.Writer.Header().Set(headerXRequestID, parsed.String())
		c.Next()
	}
}

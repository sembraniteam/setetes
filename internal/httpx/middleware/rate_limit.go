package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sembraniteam/setetes/internal/httpx/response"
)

const (
	headerRateLimit          = "X-RateLimit-Limit"
	headerRateLimitReset     = "X-RateLimit-Reset"
	headerRateLimitRemaining = "X-RateLimit-Remaining"
	headerRateLimitUsed      = "X-RateLimit-Used"
	headerRetryAfter         = "Retry-After"
)

func RateLimitByIP(t *TokenBucket) gin.HandlerFunc {
	return func(c *gin.Context) {
		metadata := t.AllowAndGetInfo(c.ClientIP())

		c.Header(headerRateLimit, fmt.Sprintf("%d", metadata.limit))
		c.Header(
			headerRateLimitRemaining,
			fmt.Sprintf("%d", metadata.remaining),
		)
		c.Header(
			headerRateLimitReset,
			fmt.Sprintf("%d", metadata.expiry.Unix()),
		)
		c.Header(headerRateLimitUsed, fmt.Sprintf("%d", metadata.used))

		if !metadata.allowed {
			retryAfter := max(
				int(metadata.expiry.Unix()-time.Now().Unix()),
				0,
			)

			c.Header(headerRetryAfter, fmt.Sprintf("%d", retryAfter))
			response.ToManyRequest(c)
			c.Abort()
			return
		}

		c.Next()
	}
}

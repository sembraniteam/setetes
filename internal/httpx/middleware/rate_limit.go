package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/megalodev/setetes/internal/httpx/response"
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
		exceeded(c, t, c.ClientIP())
	}
}

func exceeded(ctx *gin.Context, t *TokenBucket, key string) {
	allowed, remaining, expiry, limit, used := t.AllowAndGetInfo(key)

	ctx.Header(headerRateLimit, fmt.Sprintf("%d", limit))
	ctx.Header(headerRateLimitRemaining, fmt.Sprintf("%d", remaining))
	ctx.Header(headerRateLimitReset, fmt.Sprintf("%d", expiry.Unix()))
	ctx.Header(headerRateLimitUsed, fmt.Sprintf("%d", used))

	if !allowed {
		retryAfter := max(int(expiry.Unix()-time.Now().Unix()), 0)

		ctx.Header(headerRetryAfter, fmt.Sprintf("%d", retryAfter))
		ctx.Abort()
		response.ToManyRequest(ctx)

		return
	}

	ctx.Next()
}

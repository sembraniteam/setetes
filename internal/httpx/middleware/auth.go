package middleware

import (
	"log/slog"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gobwas/glob"
	"github.com/google/uuid"
	"github.com/sembraniteam/setetes/internal/cryptox/pasetox"
	"github.com/sembraniteam/setetes/internal/httpx"
	"github.com/sembraniteam/setetes/internal/httpx/response"
	"github.com/sembraniteam/setetes/internal/rbac"
)

var log = slog.Default()

const rulesLen = 3

type Config struct {
	manager  *rbac.Manager
	verifier *pasetox.Verifier
	patterns []glob.Glob
}

func NewAuthorizationConfig(
	manager *rbac.Manager,
	verifier *pasetox.Verifier,
	prefixes []string,
) *Config {
	patterns := make([]glob.Glob, 0, len(prefixes))
	for _, prefix := range prefixes {
		if pattern, err := glob.Compile(prefix); err == nil {
			patterns = append(patterns, pattern)
		} else {
			log.Error(
				"Failed to compile glob pattern",
				slog.String("prefix", prefix),
				slog.String("error", err.Error()),
			)
		}
	}

	return &Config{
		manager:  manager,
		verifier: verifier,
		patterns: patterns,
	}
}

func (c *Config) Authorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth(ctx, *c)
	}
}

func auth(
	c *gin.Context,
	config Config,
) {
	enforcer := config.manager.GetEnforcer()
	ctx := httpx.NewContext(c)
	resource := c.Request.URL.Path
	action := c.Request.Method

	for _, pattern := range config.patterns {
		if pattern.Match(c.Request.URL.Path) {
			ctx.SetUserSession(httpx.UserSession{
				ID:        uuid.Nil,
				Anonymous: true,
			})
			c.Next()
			return
		}
	}

	header := c.GetHeader("Authorization")
	if header == "" {
		log.Warn(
			"Missing Authorization header",
			slog.String("method", action),
			slog.String("url", resource),
		)
		response.Unauthorized(c)
		c.Abort()
		return
	}

	token := strings.TrimPrefix(header, "Bearer ")
	if token == header {
		log.Warn(
			"Invalid header format",
			slog.String("method", action),
			slog.String("url", resource),
		)
		response.Unauthorized(c)
		c.Abort()
		return
	}

	claims, err := config.verifier.Verify(token)
	if err != nil {
		log.Error(
			"Invalid/expired token",
			slog.String("method", action),
			slog.String("url", resource),
			slog.String("error", err.Error()),
		)
		response.Unauthorized(c)
		c.Abort()
		return
	}

	var domain string
	grouping, err := enforcer.GetFilteredGroupingPolicy(0, claims.Subject)
	if err != nil {
		log.Error(
			"Failed to get filtered grouping policy",
			slog.String("method", action),
			slog.String("url", resource),
			slog.String("error", err.Error()),
		)
	}

	for _, rule := range grouping {
		if len(rule) >= rulesLen {
			domain = rule[2]
			break
		}
	}

	if domain == "" {
		log.Warn(
			"Missing domain for subject",
			slog.String("subject", claims.Subject),
			slog.String("method", action),
			slog.String("url", resource),
		)
		response.Forbidden(c)
		c.Abort()
		return
	}

	allowed, err := enforcer.Enforce(
		claims.Subject,
		domain,
		resource,
		action,
	)
	if err != nil {
		log.Error(
			"Failed to check permissions",
			slog.String("subject", claims.Subject),
			slog.String("method", action),
			slog.String("url", resource),
			slog.String("error", err.Error()),
		)
		response.Unauthorized(c)
		c.Abort()
		return
	}

	if !allowed {
		log.Warn(
			"Permission denied",
			slog.String("subject", claims.Subject),
			slog.String("method", action),
			slog.String("url", resource),
		)
		response.Forbidden(c)
		c.Abort()
		return
	}

	parsedSubject, err := uuid.Parse(claims.Subject)
	if err != nil {
		log.Error(
			"Invalid UUID format",
			slog.String("subject", claims.Subject),
			slog.String("method", action),
			slog.String("url", resource),
			slog.String("error", err.Error()),
		)
		response.Unauthorized(c)
		c.Abort()
		return
	}
	ctx.SetUserSession(httpx.UserSession{
		ID:        parsedSubject,
		Anonymous: false,
	})
	ctx.SetUserSessionClaims(httpx.UserSessionClaims{
		Claims: *claims,
	})

	c.Next()
}

package httpx

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	requestIDContextKey         = "__RequestID__"
	userSessionContextKey       = "__UserSession__"
	userSessionClaimsContextKey = "__UserSessionClaims__"
)

type Context struct {
	ctx *gin.Context
}

func NewContext(c *gin.Context) Context {
	return Context{ctx: c}
}

func (c Context) SetRequestID(id uuid.UUID) {
	c.ctx.Set(requestIDContextKey, id.String())
}

func (c Context) GetRequestID() uuid.UUID {
	rid, ok := c.ctx.Get(requestIDContextKey)
	if !ok || rid == nil {
		return uuid.Nil
	}

	id, ok := rid.(string)
	if !ok {
		return uuid.Nil
	}

	parsed, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil
	}

	return parsed
}

func (c Context) SetUserSession(us UserSession) {
	c.ctx.Set(userSessionContextKey, us)
}

func (c Context) GetUserSession() *UserSession {
	us, ok := c.ctx.Get(userSessionContextKey)
	if !ok || us == nil {
		return nil
	}

	if u, ok := us.(UserSession); ok {
		return &([]UserSession{u})[0]
	}

	return nil
}

func (c Context) SetUserSessionClaims(us UserSessionClaims) {
	c.ctx.Set(userSessionClaimsContextKey, us)
}

func (c Context) GetUserSessionClaims() *UserSessionClaims {
	us, ok := c.ctx.Get(userSessionClaimsContextKey)
	if !ok || us == nil {
		return nil
	}

	if u, ok := us.(UserSessionClaims); ok {
		return &([]UserSessionClaims{u})[0]
	}

	return nil
}

package httpx

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIDContextKey = "__RequestID__"

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

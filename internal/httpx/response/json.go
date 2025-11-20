package response

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/megalodev/setetes/internal/httpx"
)

func BadRequest(c *gin.Context, code int16, message *Message) {
	json(c, http.StatusBadRequest, code, message, nil)
}

func Ok(c *gin.Context, message *Message, result any) {
	json(c, http.StatusOK, httpx.OKCode, message, result)
}

func ToManyRequest(c *gin.Context) {
	json(c, http.StatusTooManyRequests, httpx.TooManyRequestsCode, MsgToManyRequest, nil)
}

func InternalServerError(c *gin.Context) {
	json(c, http.StatusInternalServerError, httpx.InternalErrorCode, MsgInternalServerError, nil)
}

func UnknownError(c *gin.Context) {
	json(c, http.StatusInternalServerError, httpx.UnknownCode, MsgUnknownError, nil)
}

func Unauthorized(c *gin.Context) {
	json(c, http.StatusUnauthorized, httpx.UnauthorizedCode, MsgUnauthorized, nil)
}

func Error(c *gin.Context, errParam any) {
	switch err := errParam.(type) {
	case error:
		if err == io.EOF {
			BadRequest(c, httpx.InvalidBodyCode, MsgInvalidBody)
		}

		InternalServerError(c)
	default:
		UnknownError(c)
	}
}

func json(c *gin.Context, status int, code int16, message *Message, result any) {
	c.JSON(status, Base{Code: code, Message: message, Result: result})
}

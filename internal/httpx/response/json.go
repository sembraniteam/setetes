package response

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sembraniteam/setetes/internal/httpx"
)

func BadRequest(c *gin.Context, code int16, message *Message) {
	json(c, http.StatusBadRequest, code, message, nil)
}

func Created(c *gin.Context, message *Message, result any) {
	json(c, http.StatusCreated, httpx.OKCode, message, result)
}

func Forbidden(c *gin.Context) {
	json(
		c,
		http.StatusForbidden,
		httpx.NoPermissionCode,
		MsgNoPermission,
		nil,
	)
}

func Ok(c *gin.Context, message *Message, result any) {
	json(c, http.StatusOK, httpx.OKCode, message, result)
}

func ToManyRequest(c *gin.Context) {
	json(
		c,
		http.StatusTooManyRequests,
		httpx.TooManyRequestsCode,
		MsgToManyRequest,
		nil,
	)
}

func InternalServerError(c *gin.Context) {
	json(
		c,
		http.StatusInternalServerError,
		httpx.InternalErrorCode,
		MsgInternalServerError,
		nil,
	)
}

func UnknownError(c *gin.Context) {
	json(
		c,
		http.StatusInternalServerError,
		httpx.UnknownCode,
		MsgUnknownError,
		nil,
	)
}

func Unauthorized(c *gin.Context) {
	json(
		c,
		http.StatusUnauthorized,
		httpx.UnauthorizedCode,
		MsgUnauthorized,
		nil,
	)
}

func Error(c *gin.Context, errParam any) {
	switch err := errParam.(type) {
	case *ErrorValidate:
		BadRequest(c, *err.GetCode(), err.GetMessage())
		return
	case error:
		if err == io.EOF {
			BadRequest(c, httpx.InvalidBodyCode, MsgInvalidBody)
			return
		}

		InternalServerError(c)
		return
	default:
		UnknownError(c)
		return
	}
}

func InvalidParameter(c *gin.Context, description string) {
	BadRequest(
		c,
		httpx.InvalidParameterCode,
		NewMessage(Warning, description),
	)
}

func json(
	c *gin.Context,
	status int,
	code int16,
	message *Message,
	result any,
) {
	w := httpx.NewContext(c)
	c.JSON(
		status,
		Base{
			Code:      code,
			RequestID: w.GetRequestID().String(),
			Message:   message,
			Result:    result,
		},
	)
}

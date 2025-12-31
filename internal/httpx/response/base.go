package response

const (
	Required = "REQUIRED"
	Warning  = "WARNING"
)

var (
	MsgSuccess             = NewMessage("SUCCESS", "Success")
	MsgPong                = NewMessage("PONG", "PONG")
	MsgUnknownError        = NewMessage("UNKNOWN_ERROR", "Unknown error")
	MsgToManyRequest       = NewMessage("TO_MANY_REQUEST", "To many request")
	MsgInternalServerError = NewMessage(
		"INTERNAL_ERROR",
		"Internal server error",
	)
	MsgUnauthorized     = NewMessage("UNAUTHORIZED", "Unauthorized")
	MsgInvalidBody      = NewMessage("INVALID_BODY", "Invalid request body")
	MsgInvalidRequestID = NewMessage(
		"INVALID_X_REQUEST_ID",
		"Invalid X-Request-ID value",
	)
	MsgInvalidJSON  = NewMessage("INVALID_JSON", "Invalid JSON format")
	MsgInvalidForm  = NewMessage("INVALID_FORM", "Invalid form format")
	MsgInvalidQuery = NewMessage("INVALID_QUERY", "Invalid query format")
	MsgNoPermission = NewMessage(
		"NO_PERMISSION",
		"You don't have access to this resource",
	)
)

type (
	Message struct {
		Key         string `json:"key"`
		Description string `json:"description"`
	}

	Base struct {
		Code      int16    `json:"code"`
		RequestID string   `json:"request_id"`
		Message   *Message `json:"message"`
		Result    any      `json:"result"`
	}
)

func NewMessage(key, description string) *Message {
	return &Message{
		Key:         key,
		Description: description,
	}
}

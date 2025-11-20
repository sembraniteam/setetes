package response

var (
	MsgUnknownError        = NewMessage("UNKNOWN_ERROR", "Unknown error", nil)
	MsgToManyRequest       = NewMessage("TO_MANY_REQUEST", "To many request", nil)
	MsgInternalServerError = NewMessage("INTERNAL_ERROR", "Internal server error", nil)
	MsgPong                = NewMessage("PONG", "PONG", nil)
	MsgUnauthorized        = NewMessage("UNAUTHORIZED", "Unauthorized", nil)
	MsgInvalidBody         = NewMessage("INVALID_BODY", "Invalid request body", nil)
	MsgInvalidRequestID    = NewMessage("INVALID_X_REQUEST_ID", "Invalid header value", nil)
)

type (
	Message struct {
		Key         string  `json:"key"`
		Description string  `json:"description"`
		Field       *string `json:"field"`
	}

	Base struct {
		Code      int16    `json:"code"`
		RequestID string   `json:"request_id"`
		Message   *Message `json:"message"`
		Result    any      `json:"result"`
	}

	BaseEntries[T any] struct {
		Entries       []T   `json:"entries"`
		HasReachedMax bool  `json:"has_reached_max"`
		TotalPages    int64 `json:"total_pages"`
	}
)

func NewMessage(key, description string, field *string) *Message {
	return &Message{
		Key:         key,
		Description: description,
		Field:       field,
	}
}

func (m *Message) WithField(s string) *Message {
	cp := *m
	cp.Field = &s
	return &cp
}

func Entries[T any](entries []T, hasReachedMax bool, totalPages int64) *BaseEntries[T] {
	return &BaseEntries[T]{
		Entries:       entries,
		HasReachedMax: hasReachedMax,
		TotalPages:    totalPages,
	}
}

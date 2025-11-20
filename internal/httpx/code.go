package httpx

const (
	OKCode                int16 = 0
	UnknownCode           int16 = 1000
	UnauthorizedCode      int16 = 1001
	TooManyRequestsCode   int16 = 1002
	InternalErrorCode     int16 = 1003
	InvalidBodyCode       int16 = 1100
	RequiredKeyCode       int16 = 1101
	InvalidJSONCode       int16 = 1102
	InvalidFormCode       int16 = 1103
	InvalidFileCode       int16 = 1200
	InvalidValueCode      int16 = 1201
	InvalidQueryCode      int16 = 1300
	RequiredQueryCode     int16 = 1301
	InvalidCredentialCode int16 = 2000
	NoPermissionCode      int16 = 2001
)

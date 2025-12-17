package response

import (
	"errors"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entrans "github.com/go-playground/validator/v10/translations/en"
	"github.com/sembraniteam/setetes/internal/httpx"
)

const n = 2

var (
	trans    ut.Translator
	validate = validator.New()
)

type ErrorValidate struct {
	code    *int16
	message *Message
}

func (v *ErrorValidate) GetCode() *int16 {
	return v.code
}

func (v *ErrorValidate) GetMessage() *Message {
	return v.message
}

func init() {
	translator, _ := ut.New(en.New()).GetTranslator("en")
	trans = translator

	if err := entrans.RegisterDefaultTranslations(validate, trans); err != nil {
		panic(err)
	}

	if err := validate.RegisterValidation("password", validatePassword); err != nil {
		panic(err)
	}

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("form"), ",", n)[0]
		if name == "" || name == "-" {
			name = strings.SplitN(fld.Tag.Get("json"), ",", n)[0]
		}

		if name == "-" {
			return ""
		}

		return name
	})
}

func ValidateJSON[T any](c *gin.Context) (*T, *ErrorValidate) {
	body := new(T)
	if err := c.ShouldBindJSON(body); err != nil {
		return nil, &ErrorValidate{
			code:    &httpx.InvalidJSONCode,
			message: MsgInvalidJSON,
		}
	}

	if err := validateStruct(body); err != nil {
		return nil, err
	}

	return body, nil
}

func ValidateForm[T any](c *gin.Context) (*T, *ErrorValidate) {
	body := new(T)
	if err := c.ShouldBindWith(body, binding.FormMultipart); err != nil {
		return nil, &ErrorValidate{
			code:    &httpx.InvalidFormCode,
			message: MsgInvalidForm,
		}
	}

	if err := validateStruct(body); err != nil {
		return nil, err
	}

	return body, nil
}

func ValidateQuery[T any](c *gin.Context) (*T, *ErrorValidate) {
	body, err := ValidateForm[T](c)
	if err != nil {
		return nil, &ErrorValidate{
			code:    &httpx.InvalidQueryCode,
			message: MsgInvalidQuery,
		}
	}

	return body, nil
}

func validateStruct(data any) *ErrorValidate {
	if err := validate.Struct(data); err != nil {
		var e validator.ValidationErrors
		if errors.As(err, &e) {
			firstErr := e[0]
			errMsg := firstErr.Translate(trans)

			dataType := reflect.TypeOf(data)
			if dataType.Kind() == reflect.Ptr {
				dataType = dataType.Elem()
			}

			if field, found := dataType.FieldByName(firstErr.StructField()); found {
				if customMsg := reason(field, firstErr.Tag()); customMsg != "" {
					errMsg = customMsg
				}
			}

			return &ErrorValidate{
				code:    &httpx.InvalidValueCode,
				message: NewMessage(Required, errMsg),
			}
		}

		return &ErrorValidate{
			code:    &httpx.InvalidBodyCode,
			message: MsgInvalidBody,
		}
	}

	return nil
}

func validatePassword(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() != reflect.String {
		return false
	}

	password := field.String()
	if password == "" {
		return false
	}

	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSymbol := false
	specialChars := "@#!$%^&*()-_=+?><~.,"

	for _, char := range password {
		if char >= 'A' && char <= 'Z' {
			hasUpper = true
		}

		if char >= 'a' && char <= 'z' {
			hasLower = true
		}

		if char >= '0' && char <= '9' {
			hasNumber = true
		}

		if strings.ContainsRune(specialChars, char) {
			hasSymbol = true
		}

		if hasUpper && hasLower && hasNumber && hasSymbol {
			return true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSymbol
}

func reason(field reflect.StructField, tag string) string {
	reasonTag := field.Tag.Get("reason")
	rules := strings.SplitSeq(reasonTag, ";")

	for rule := range rules {
		parts := strings.SplitN(rule, "=", n)
		if len(parts) == 2 && parts[0] == tag {
			return parts[1]
		}
	}

	return ""
}

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
	"github.com/megalodev/setetes/internal/httpx"
)

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

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
		if name == "" || name == "-" {
			name = strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		}

		if name == "-" {
			return ""
		}

		return name
	})
}

func ValidateJson[T any](c *gin.Context) (*T, *int16, *Message) {
	body := new(T)
	if err := c.ShouldBindJSON(body); err != nil {
		return nil, &httpx.InvalidJSONCode, MsgInvalidJSON
	}

	code, err := validateStruct(body)
	if err != nil || code == nil {
		return nil, code, err
	}

	return body, nil, nil
}

func ValidateForm[T any](c *gin.Context) (*T, *int16, *Message) {
	body := new(T)
	if err := c.ShouldBindWith(body, binding.FormMultipart); err != nil {
		return nil, &httpx.InvalidFormCode, MsgInvalidForm
	}

	code, err := validateStruct(body)
	if err != nil || code != nil {
		return nil, code, err
	}

	return body, nil, nil
}

func ValidateQuery[T any](c *gin.Context) (*T, *int16, *Message) {
	body, code, err := ValidateForm[T](c)

	if err != nil || code != nil {
		return nil, &httpx.InvalidQueryCode, MsgInvalidQuery
	}

	return body, nil, nil
}

func validateStruct(data any) (*int16, *Message) {
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
				if customMsg := customReason(field, firstErr.Tag()); customMsg != "" {
					errMsg = customMsg
				}
			}

			return &httpx.InvalidValueCode, NewMessage("REQUIRED", errMsg)
		}

		return &httpx.InvalidBodyCode, MsgInvalidBody
	}

	return nil, nil
}

func customReason(field reflect.StructField, tag string) string {
	reasonTag := field.Tag.Get("reason")
	rules := strings.Split(reasonTag, ";")

	for _, rule := range rules {
		parts := strings.SplitN(rule, "=", 2)
		if len(parts) == 2 && parts[0] == tag {
			return parts[1]
		}
	}

	return ""
}

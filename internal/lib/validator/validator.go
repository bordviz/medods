package validator

import (
	"fmt"
	"medods/internal/lib/customerror"
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	once     sync.Once
)

func getValidator() *validator.Validate {
	once.Do(func() {
		validate = validator.New()
	})
	return validate
}

func Validate(model interface{}) customerror.CustomError {
	var errMsgs []string

	validate := getValidator()
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("json")
	})

	err := validate.Struct(model)

	if err != nil {
		validErr := err.(validator.ValidationErrors)
		for _, errMsg := range validErr {
			switch errMsg.ActualTag() {
			case "required":
				errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required", errMsg.Field()))
			case "email":
				errMsgs = append(errMsgs, fmt.Sprintf("field %s must be a valid email", errMsg.Field()))
			default:
				errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", errMsg.Field()))
			}
		}
		return customerror.NewCustomError(fmt.Sprintf("validation error: %s", strings.Join(errMsgs, ", ")), 422)
	}
	return nil
}

package utils

import (
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterTagNameFunc(func(fid reflect.StructField) string {
		name := strings.SplitN(fid.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			name = ""
		}
		return name
	})
	validate.RegisterValidation("futureDate", futureDate)
}

func ValidateStruct(input interface{}) error {
	return validate.Struct(input)
}

func futureDate(fl validator.FieldLevel) bool {
	fieldValue := fl.Field().Interface()
	if fieldValue == nil {
		return true
	}

	t, ok := fieldValue.(*time.Time)
	if !ok {
		return false
	}

	return t.After(time.Now())
}

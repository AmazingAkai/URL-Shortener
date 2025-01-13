package utils

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	validate           *validator.Validate
	alphanumericRegexp = regexp.MustCompile("^[a-zA-Z0-9]+$")
)

func init() {
	validate = validator.New()
	validate.RegisterValidation("alphanumeric", alphanumeric)
	validate.RegisterTagNameFunc(func(fid reflect.StructField) string {
		name := strings.SplitN(fid.Tag.Get("schema"), ",", 2)[0]
		if name == "-" {
			name = ""
		}
		return name
	})
}

func ValidateStruct(input interface{}) error {
	return validate.Struct(input)
}

func alphanumeric(fl validator.FieldLevel) bool {
	return alphanumericRegexp.MatchString(fl.Field().String())
}

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
	validate.RegisterValidation("futureDate", futureDate)
	validate.RegisterTagNameFunc(func(fid reflect.StructField) string {
		name := strings.SplitN(fid.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			name = ""
		}
		return name
	})
}

func ValidateStruct(input interface{}) error {
	return validate.Struct(input)
}

func futureDate(fl validator.FieldLevel) bool {
	if timeVal, ok := fl.Field().Interface().(time.Time); ok {
		currentTime := time.Now()
		twentyEightDaysFromNow := currentTime.Add(28 * 24 * time.Hour)

		if timeVal.After(currentTime) && timeVal.Before(twentyEightDaysFromNow) {
			return true
		}
	}

	return false
}

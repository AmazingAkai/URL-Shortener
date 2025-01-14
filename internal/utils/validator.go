package utils

import (
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	validate           *validator.Validate
	alphanumericRegexp = regexp.MustCompile("^[a-zA-Z0-9]+$")
	WEB_URL            = os.Getenv("WEB_URL")
	RESERVED_URLS      = []string{
		"static",
		"urls",
		"login",
		"register",
		"logout",
	}
)

func init() {
	validate = validator.New()
	validate.RegisterValidation("validShortUrl", validShortUrl)
	validate.RegisterValidation("validLongUrl", validLongUrl)
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

func validShortUrl(fl validator.FieldLevel) bool {
	return !SliceContains(RESERVED_URLS, fl.Field().String())
}

func validLongUrl(fl validator.FieldLevel) bool {
	url := fl.Field().String()
	return !strings.Contains(url, WEB_URL) &&
		(strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://"))
}

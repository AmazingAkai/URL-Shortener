package utils

import (
	"fmt"
	"net/http"

	"github.com/AmazingAkai/URL-Shortener/app/internal/log"

	"github.com/go-playground/validator/v10"
)

type ErrorMap map[string][]string

func ValidationError(w http.ResponseWriter, vErr error) {
	resp := ErrorMap{}

	switch err := vErr.(type) {
	case validator.ValidationErrors:
		for _, e := range err {
			field := e.Field()
			msg := checkTagRules(e)
			resp[field] = append(resp[field], msg)
		}
	default:
		resp["non_field_error"] = append(resp["non_field_error"], err.Error())
	}
	ErrorResponse(w, http.StatusUnprocessableEntity, resp)
}

func BadRequestError(w http.ResponseWriter) {
	ErrorResponse(w, http.StatusUnprocessableEntity, "unable to process request")
}

func NotFoundError(w http.ResponseWriter) {
	ErrorResponse(w, http.StatusNotFound, "not found")
}

func ServerError(w http.ResponseWriter, err error) {
	log.Errorf("Internal server error: %v", err)
	ErrorResponse(w, http.StatusInternalServerError, "internal error")
}

func ErrorResponse(w http.ResponseWriter, code int, errs interface{}) {
	WriteJSON(w, code, Map{"errors": errs})
}

func checkTagRules(e validator.FieldError) (errMsg string) {
	tag, field, param, value := e.ActualTag(), e.Field(), e.Param(), e.Value()

	switch tag {
	case "required":
		errMsg = "this field is required"
	case "email":
		errMsg = fmt.Sprintf("%q is not a valid email", value)
	case "min":
		errMsg = fmt.Sprintf("%s must be greater than %v", field, param)
	case "max":
		errMsg = fmt.Sprintf("%s must be less than %v", field, param)
	default:
		errMsg = fmt.Sprintf("%s is not valid", field)
	}

	return
}

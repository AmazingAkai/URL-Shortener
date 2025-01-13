package utils

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AmazingAkai/URL-Shortener/internal/views/partials"
	"github.com/go-playground/validator/v10"
)

func ValidationError(w http.ResponseWriter, r *http.Request, validationError error) {
	var resp []string

	switch err := validationError.(type) {
	case validator.ValidationErrors:
		resp = make([]string, 0, len(err))
		for _, e := range err {
			resp = append(resp, checkTagRules(e))
		}
	default:
		resp = []string{err.Error()}
	}
	ErrorResponse(w, r, http.StatusBadRequest, resp)
}

func ParseFormError(w http.ResponseWriter, r *http.Request, err error) {
	ErrorResponse(w, r, http.StatusBadRequest, []string{err.Error()})
}

func ServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Internal server error: %v", err)
	ErrorResponse(w, r, http.StatusInternalServerError, []string{"An internal server error occurred."})
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, code int, errs []string) {
	w.WriteHeader(code)
	partials.Error(errs).Render(r.Context(), w)
}

func checkTagRules(e validator.FieldError) (errMsg string) {
	tag, field, param, value := e.ActualTag(), e.Field(), e.Param(), e.Value()

	switch tag {
	case "required":
		errMsg = fmt.Sprintf("%s is required.", field)
	case "email":
		errMsg = fmt.Sprintf("%q is not a valid email.", value)
	case "min":
		errMsg = fmt.Sprintf("%s length must be greater than %v.", field, param)
	case "max":
		errMsg = fmt.Sprintf("%s must be less than %v.", field, param)
	case "alphanumeric":
		errMsg = fmt.Sprintf("%s must be alphanumeric.", field)
	case "url":
		errMsg = fmt.Sprintf("%q is not a valid URL.", value)
	default:
		errMsg = fmt.Sprintf("%s is not valid.", field)
	}

	return
}

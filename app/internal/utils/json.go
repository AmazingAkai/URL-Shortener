package utils

import (
	"io"
	"net/http"

	jsoniter "github.com/json-iterator/go"

	"github.com/AmazingAkai/URL-Shortener/app/internal/log"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Map map[string]interface{}

func WriteJSON(w http.ResponseWriter, code int, data interface{}) {
	jsonBytes, err := json.Marshal(data)

	if err != nil {
		ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(jsonBytes)

	if err != nil {
		log.Errorf("Error writing response: %v", err)
	}
}

func ReadJSON(body io.Reader, input interface{}) error {
	return json.NewDecoder(body).Decode(input)
}

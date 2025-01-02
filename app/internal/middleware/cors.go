package middleware

import (
	"net/http"
)

const (
	ALLOW_ORIGIN      = "*"
	ALLOW_METHODS     = "GET, POST, OPTIONS"
	ALLOW_HEADERS     = "Accept, Authorization"
	ALLOW_CREDENTIALS = "false"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", ALLOW_ORIGIN)
		w.Header().Set("Access-Control-Allow-Methods", ALLOW_METHODS)
		w.Header().Set("Access-Control-Allow-Headers", ALLOW_HEADERS)
		w.Header().Set("Access-Control-Allow-Credentials", ALLOW_CREDENTIALS)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

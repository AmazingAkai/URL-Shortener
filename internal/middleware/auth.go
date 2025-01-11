package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/AmazingAkai/URL-Shortener/internal/store"
	"github.com/AmazingAkai/URL-Shortener/internal/utils"
	"github.com/AmazingAkai/URL-Shortener/internal/utils/constants"
)

func Auth(s *store.Storage) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionToken, err := r.Cookie("session_token")
			if err != nil || sessionToken.Value == "" {
				next.ServeHTTP(w, r)
				return
			}
			csrfToken, err := r.Cookie("csrf_token")
			if err != nil || csrfToken.Value == "" {
				utils.ErrorResponse(w, http.StatusUnauthorized, "csrf token not found")
				return
			}

			if csrfToken.Value != r.Header.Get("X-CSRF-TOKEN") {
				fmt.Printf("CSRFTOKEN COOKIE VALUE: %v\n", csrfToken.Value)
				fmt.Printf("X-CSRF-TOKEN HEADER VALUE: %v\n", r.Header.Get("X-CSRF-TOKEN"))
				utils.ErrorResponse(w, http.StatusUnauthorized, "invalid csrf token")
				return
			}

			session, err := s.Sessions.Get(sessionToken.Value)
			if err != nil {
				utils.ErrorResponse(w, http.StatusUnauthorized, "invalid session token")
				return
			}

			ctx := context.WithValue(r.Context(), constants.SESSION_KEY, session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

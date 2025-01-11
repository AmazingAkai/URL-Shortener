package middleware

import (
	"context"
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

			session, err := s.Sessions.Get(sessionToken.Value)
			if err != nil {
				utils.ErrorResponse(w, r, http.StatusUnauthorized, []string{"Invalid session token."})
				return
			}

			ctx := context.WithValue(r.Context(), constants.SESSION_KEY, session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

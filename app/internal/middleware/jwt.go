package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/AmazingAkai/URL-Shortener/app/internal/models"
	"github.com/AmazingAkai/URL-Shortener/app/internal/utils"
	"github.com/AmazingAkai/URL-Shortener/app/internal/utils/constants"

	"github.com/golang-jwt/jwt/v4"
)

func JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := utils.ValidateJWT(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
		if err != nil || !token.Valid {
			next.ServeHTTP(w, r)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["email"] == nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), constants.USER_KEY, &models.UserOut{
			ID:    int(claims["id"].(float64)),
			Email: claims["email"].(string),
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

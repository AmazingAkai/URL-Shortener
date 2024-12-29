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

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		tokenString := strings.Split(authHeader, "Bearer ")[1]
		if tokenString == "" {
			next.ServeHTTP(w, r)
			return
		}

		token, err := utils.ValidateJWT(tokenString)
		if err != nil || !token.Valid {
			next.ServeHTTP(w, r)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["email"] == nil {
			next.ServeHTTP(w, r)
			return
		}

    // TODO: check wheter the user exists in the database

		user := &models.User{
			ID:    int(claims["id"].(float64)),
			Email: claims["email"].(string),
		}

		ctx := context.WithValue(r.Context(), constants.USER_KEY, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

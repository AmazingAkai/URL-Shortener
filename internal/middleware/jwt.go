package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/AmazingAkai/URL-Shortener/internal/store"
	"github.com/AmazingAkai/URL-Shortener/internal/utils/constants"

	"github.com/golang-jwt/jwt/v4"
)

var JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))

func JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := validateJWT(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
		if err != nil || !token.Valid {
			next.ServeHTTP(w, r)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["email"] == nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), constants.USER_KEY, &store.User{
			ID:    int(claims["id"].(float64)),
			Email: claims["email"].(string),
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GenerateJWT(user *store.User) (string, int64, error) {
	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(JWT_SECRET)
	if err != nil {
		return "", 0, fmt.Errorf("failed to sign the token: %v", err)
	}

	return token, claims["exp"].(int64), nil
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return JWT_SECRET, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

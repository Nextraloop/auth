package middleware

import (
	"auth-service/handlers"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Extract token from Authorization header
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse the token
		claims := &handlers.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(handlers.JwtKey), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

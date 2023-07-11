package middleware

import (
	"api/config"
	"api/helper"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If token is not found in the cookie, check the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" { // If Authorization header is not found also
			helper.ResponseJSON(w, http.StatusUnauthorized, "error", "Unauthorized", nil)
			return
		}

		tokenString := authHeader

		claims := &config.JWTClaim{}

		// Parse tokens
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return config.JWT_KEY, nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
				// Token invalid
				helper.ResponseJSON(w, http.StatusUnauthorized, "error", "Unauthorized", nil)
				return
			} else if errors.Is(err, jwt.ErrTokenExpired) {
				// Token expired
				helper.ResponseJSON(w, http.StatusUnauthorized, "error", "Unauthorized, token expired!", nil)
				return
			} else {
				helper.ResponseJSON(w, http.StatusUnauthorized, "error", "Couldn't handle this token", nil)
				return
			}
		}

		if !token.Valid {
			helper.ResponseJSON(w, http.StatusUnauthorized, "error", "Unauthorized", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}

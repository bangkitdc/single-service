package selfcontroller

import (
	"api/config"
	"api/helper"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func GetSelf(w http.ResponseWriter, r *http.Request) {
	// Get the current user information
	// Check the Authorization header
	tokenString := r.Header.Get("Authorization")
	tokenString = tokenString[len("Bearer "):] // Remove the "Bearer " prefix

	// Parse the token and extract the claims
	claims := &config.JWTClaim{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return config.JWT_KEY, nil
	})

	if err != nil {
		helper.ResponseJSON(w, http.StatusUnauthorized, "error", "Unauthorized, Login first!", nil)
		return
	}

	// Create the response data
	responseData := struct {
		Username string `json:"username"`
		Name     string `json:"name"`
	}{
		Username: claims.Username,
		Name:     claims.Name,
	}

	helper.ResponseJSON(w, http.StatusOK, "success", "User retrieved successfully", responseData)
}

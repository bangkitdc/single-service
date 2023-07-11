package authcontroller

import (
	"api/config"
	"api/helper"
	"api/models"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func Login(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var requestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		helper.ResponseJSON(w, http.StatusBadRequest, "error", "Invalid request body", nil)
		return
	}

	// Validate the required fields
	if helper.IsEmpty(requestBody.Username) || helper.IsEmpty(requestBody.Password) {
		helper.ResponseJSON(w, http.StatusBadRequest, "error", "Missing required fields", nil)
		return
	}

	// Get user
	var user models.User

	if err := models.DB.Where("username = ?", requestBody.Username).First(&user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			helper.ResponseJSON(w, http.StatusUnauthorized, "error", "Invalid credentials", nil)
			return

		default:
			helper.ResponseJSON(w, http.StatusInternalServerError, "error", err.Error(), nil)
			return
		}
	}

	// Check incredentials
	if user.Password != requestBody.Password {
		helper.ResponseJSON(w, http.StatusUnauthorized, "error", "Invalid credentials", nil)
		return
	}

	// JWT token
	expiredTime := time.Now().Add(time.Hour * 1)
	claims := &config.JWTClaim{
		Username: user.Username,
		Name:     user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-jwt-mux",
			ExpiresAt: jwt.NewNumericDate(expiredTime),
		},
	}

	tokenAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := tokenAlgo.SignedString(config.JWT_KEY)
	if err != nil {
		helper.ResponseJSON(w, http.StatusInternalServerError, "error", err.Error(), nil)
		return
	}

	// Create a custom struct for the response data
	type UserData struct {
		Username string `json:"username"`
		Name     string `json:"name"`
	}

	// Create the final response
	responseData := struct {
		User  *UserData `json:"user"`
		Token string    `json:"token"`
	}{
		User: &UserData{
			Username: user.Username,
			Name:     user.Name,
		},
		Token: token,
	}

	// Set token to response
	helper.ResponseJSON(w, http.StatusOK, "success", "Login successful", responseData)
}

package helper

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ResponseJSON(w http.ResponseWriter, code int, status string, message string, data interface{}) {
	response := APIResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}

	responseBytes, _ := json.Marshal(response)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(responseBytes)
}

func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

func IncrementID(id string) string {
	// Convert the ID to an integer
	num, err := strconv.Atoi(id)
	if err != nil {
		// Handle the error when the ID cannot be converted to an integer

		return "error"
	}

	// Increment the numeric value
	num++

	// Convert the incremented value back to a string
	incrementedID := strconv.Itoa(num)

	return incrementedID
}

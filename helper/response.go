package helper

import (
	"encoding/json"
	"net/http"
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

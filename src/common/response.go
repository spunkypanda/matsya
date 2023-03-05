package common

import (
	"encoding/json"
	"net/http"
)

func SendSuccessResponse(w http.ResponseWriter, data any, message string) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]any{}

	if message != "" {
		response["message"] = message
	}

	if data != nil {
		response["data"] = data
	}

	json.NewEncoder(w).Encode(response)
}

func SendErrorResponse(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")

	if status == 0 {
		status = http.StatusInternalServerError
	}

	w.WriteHeader(status)

	response := map[string]interface{}{"message": message}
	json.NewEncoder(w).Encode(response)
}

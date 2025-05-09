package api

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func responseWithJSON(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(errorResponse{
		Code:    statusCode,
		Message: message,
	})
}

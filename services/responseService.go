package services

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func ResponseWithMessage(w http.ResponseWriter, statusCode int, headers map[string]string, message string, requestID string) {

	for k, v := range headers {
		w.Header().Set(k, v)
	}

	w.Header().Set("Content-Type", "application/json")
	if requestID != "" {
		w.Header().Set("X-Request-ID", requestID)
	}
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(APIResponse{
		Message: message,
	})
}

func ResponseWithData(w http.ResponseWriter, statusCode int, headers map[string]string, message string, data interface{}, requestID string) {

	for k, v := range headers {
		w.Header().Set(k, v)
	}

	w.Header().Set("Content-Type", "application/json")
	if requestID != "" {
		w.Header().Set("X-Request-ID", requestID)
	}
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(APIResponse{
		Message: message,
		Data:    data,
	})
}

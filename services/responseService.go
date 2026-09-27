package services

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func responseStatus(statusCode int) string {
	switch {
	case statusCode == http.StatusAccepted:
		return "accepted"
	case statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError:
		return "failed"
	case statusCode >= http.StatusInternalServerError:
		return "error"
	default:
		return "success"
	}
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
		Status:  responseStatus(statusCode),
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
		Status:  responseStatus(statusCode),
		Message: message,
		Data:    data,
	})
}

func Redirect(w http.ResponseWriter, statusCode int, headers map[string]string, url string, r *http.Request) {
	for k, v := range headers {
		w.Header().Set(k, v)
	}
	http.Redirect(w, r, url, statusCode)
}

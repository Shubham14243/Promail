package models

type LogData struct {
	RequestID    string `json:"request_id,omitempty"`
	Endpoint     string `json:"endpoint,omitempty"`
	Method       string `json:"method,omitempty"`
	Operation    string `json:"operation,omitempty"`
	Status       string `json:"status,omitempty"`
	UserID       string `json:"user_id,omitempty"`
	ResourceID   string `json:"resource_id,omitempty"`
	ResponseCode int    `json:"response_code,omitempty"`
	Message      string `json:"message,omitempty"`
	Error        string `json:"error,omitempty"`
}

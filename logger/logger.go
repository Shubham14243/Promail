package logger

import (
	"log/slog"
	"os"
	"promail/models"
)

var Log *slog.Logger

func Init() {
	Log = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
	)
}

func Error(logData models.LogData) {
	Log.Error(
		logData.Message,
		"request_id", logData.RequestID,
		"endpoint", logData.Endpoint,
		"method", logData.Method,
		"operation", logData.Operation,
		"status", logData.Status,
		"user_id", logData.UserID,
		"resource_id", logData.ResourceID,
		"response_code", logData.ResponseCode,
		"error", logData.Error,
	)
}

func Debug(logData models.LogData) {
	Log.Debug(
		logData.Message,
		"request_id", logData.RequestID,
		"endpoint", logData.Endpoint,
		"method", logData.Method,
		"operation", logData.Operation,
		"status", logData.Status,
		"user_id", logData.UserID,
		"resource_id", logData.ResourceID,
		"response_code", logData.ResponseCode,
		"error", logData.Error,
	)
}

func Info(logData models.LogData) {
	Log.Info(
		logData.Message,
		"request_id", logData.RequestID,
		"endpoint", logData.Endpoint,
		"method", logData.Method,
		"operation", logData.Operation,
		"status", logData.Status,
		"user_id", logData.UserID,
		"resource_id", logData.ResourceID,
		"response_code", logData.ResponseCode,
		"error", logData.Error,
	)
}

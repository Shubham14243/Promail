package middlewares

import (
	"context"
	"net/http"
	"strconv"

	"promail/logger"
	"promail/models"
	"promail/services"
	"strings"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func Auth(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		logdata := models.LogData{
			RequestID: r.Context().Value(RequestIDKey).(string),
			Endpoint:  r.RequestURI,
			Method:    r.Method,
			Operation: "User Authorization",
			Status:    "Init",
			UserID:    "",
			Message:   "Authorization started.",
		}
		logger.Info(logdata)

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			logdata.Message = "Missing auth header."
			logdata.Status = "Failure"
			logdata.ResponseCode = http.StatusUnauthorized
			logdata.Error = ""
			logger.Info(logdata)
			services.ResponseWithMessage(
				w,
				http.StatusUnauthorized,
				nil,
				"Missing auth header.",
				logdata.RequestID,
			)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := services.ValidateAccessToken(tokenString)

		if err != nil {
			logdata.Message = "Authorization failure."
			logdata.Status = "Failure"
			logdata.ResponseCode = http.StatusUnauthorized
			logdata.Error = err.Error()
			logger.Info(logdata)
			services.ResponseWithMessage(
				w,
				http.StatusUnauthorized,
				nil,
				err.Error(),
				logdata.RequestID,
			)
			return
		}

		userID := int64(claims["user_id"].(float64))

		ctx := context.WithValue(
			r.Context(),
			UserIDKey,
			userID,
		)

		logdata.Message = "Autorization successful."
		logdata.Status = "Success"
		logdata.ResponseCode = 0
		logdata.Error = ""
		logdata.UserID = strconv.FormatInt(userID, 10)
		logger.Info(logdata)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

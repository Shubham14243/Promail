package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"promail/logger"
	"promail/middlewares"
	"promail/models"
	"promail/repositories"
	"promail/services"
)

type EmailHandler struct {
	EmailRepo     *repositories.EmailRepository
	AppConfigRepo *repositories.AppConfigRepository
}

func (h *EmailHandler) SendEmail(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middlewares.UserIDKey).(int64)

	logdata := models.LogData{
		RequestID: r.Context().Value(middlewares.RequestIDKey).(string),
		Endpoint:  r.RequestURI,
		Method:    r.Method,
		Operation: "Email Send Request",
		Status:    "Init",
		UserID:    strconv.FormatInt(userID, 10),
		Message:   "Email sending initiated.",
	}
	logger.Info(logdata)

	var req models.EmailSend

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logdata.Message = "Request body parsing failed."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, "Invalid request body. 'app_id', 'to', 'subject', 'body' required.", logdata.RequestID)
		return
	}

	exists, err := h.AppConfigRepo.AppConfigExistsByAppID(int64(req.AppID), userID)
	if !exists {
		logdata.Message = "App config not found."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusNotFound
		logdata.Error = ""
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusNotFound, nil, "App config not found.", logdata.RequestID)
		return
	}

	appConf, err := h.AppConfigRepo.GetAppConfigs(int64(req.AppID), userID)
	if err != nil {
		logdata.Message = "No app config data found."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusNotFound
		logdata.Error = err.Error()
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusNotFound, nil, "No app config data found.", logdata.RequestID)
		return
	}

	decrypted_password, err := services.Decrypt(appConf.SMTPPassword)
	if err != nil {
		logdata.Message = "Password decryption failed."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusNotFound, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	if err := services.SendEmail(appConf.SMTPHost, appConf.SMTPPort, appConf.SMTPUsername, decrypted_password, req.To, req.Subject, req.Body); err != nil {
		logdata.Message = "Email sending failure"
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	logdata.Message = "Email sending successful."
	logdata.Status = "Success"
	logdata.ResponseCode = http.StatusOK
	logdata.Error = ""
	logger.Info(logdata)
	services.ResponseWithMessage(w, http.StatusCreated, nil, "Email sent successfully.", logdata.RequestID)

}

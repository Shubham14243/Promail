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

type AppConfigHandler struct {
	AppConfigRepo *repositories.AppConfigRepository
}

func (h *AppConfigHandler) CreateAppConfig(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middlewares.UserIDKey).(int64)

	logdata := models.LogData{
		RequestID: r.Context().Value(middlewares.RequestIDKey).(string),
		Endpoint:  r.RequestURI,
		Method:    r.Method,
		Operation: "App Config Create",
		Status:    "Init",
		UserID:    strconv.FormatInt(userID, 10),
		Message:   "App Config Creation initiated.",
	}
	logger.Info(logdata)

	var req models.AppConfigCreate

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logdata.Message = "Request body parsing failed."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, "Invalid request body. 'app_id', 'smtp_host', 'smtp_port', 'smtp_username', 'smtp_password', 'open_track', 'click_track' required.", logdata.RequestID)
		return
	}

	if err := services.ValidateAppConfigCreate(req); err != nil {
		logdata.Message = "Request body validation failed."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, err.Error(), logdata.RequestID)
		return
	}

	exists, err := h.AppConfigRepo.AppConfigExistsByAppID(req.AppID, userID)
	if err != nil {
		logdata.Message = "App config exist check failure"
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	if exists {
		logdata.Message = "App config already exists."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = ""
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, "App config already exists.", logdata.RequestID)
		return
	}

	encryptedPassword, err := services.Encrypt(req.SMTPPassword)
	if err != nil {
		logdata.Message = "Password encryption failure"
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	appConfig := models.AppConfigCreate{
		AppID:        req.AppID,
		SMTPHost:     req.SMTPHost,
		SMTPPort:     req.SMTPPort,
		SMTPUsername: req.SMTPUsername,
		SMTPPassword: encryptedPassword,
		OpenTrack:    req.OpenTrack,
		ClickTrack:   req.ClickTrack,
	}

	if err := h.AppConfigRepo.CreateAppConfig(appConfig); err != nil {
		logdata.Message = "App config creation failure"
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	logdata.Message = "App config creation successful."
	logdata.Status = "Success"
	logdata.ResponseCode = http.StatusOK
	logdata.Error = ""
	logger.Info(logdata)
	services.ResponseWithMessage(w, http.StatusCreated, nil, "App config created successfully.", logdata.RequestID)
}

func (h *AppConfigHandler) GetAppConfigData(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middlewares.UserIDKey).(int64)
	idStr := r.PathValue("appID")
	appID, _ := strconv.Atoi(idStr)

	logdata := models.LogData{
		RequestID:  r.Context().Value(middlewares.RequestIDKey).(string),
		Endpoint:   r.RequestURI,
		Method:     r.Method,
		Operation:  "Fetch App Config Data",
		Status:     "Init",
		UserID:     strconv.FormatInt(userID, 10),
		Message:    "User app configs fetch initiated.",
		ResourceID: strconv.Itoa(appID),
	}
	logger.Info(logdata)

	exists, err := h.AppConfigRepo.AppConfigExistsByAppID(int64(appID), userID)
	if !exists {
		logdata.Message = "App config not found."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusNotFound
		logdata.Error = ""
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusNotFound, nil, "App config not found.", logdata.RequestID)
		return
	}

	appConfigData, err := h.AppConfigRepo.GetAppConfigs(int64(appID), userID)
	if err != nil {
		logdata.Message = "No app config data found."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusNotFound
		logdata.Error = err.Error()
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusNotFound, nil, "No app config data found.", logdata.RequestID)
		return
	}

	appConfigData.SMTPPassword, err = services.Decrypt(appConfigData.SMTPPassword)
	if err != nil {
		logdata.Message = "Password decryption failure"
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	logdata.Message = "App config data fetch successful."
	logdata.Status = "Success"
	logdata.ResponseCode = http.StatusOK
	logdata.Error = ""
	logger.Info(logdata)
	services.ResponseWithData(w, http.StatusOK, nil, "App config found successfully", appConfigData, logdata.RequestID)

}

func (h *AppConfigHandler) UpdateAppConfig(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middlewares.UserIDKey).(int64)
	idStr := r.PathValue("appID")
	appID, _ := strconv.Atoi(idStr)

	var req models.AppConfigUpdate

	logdata := models.LogData{
		RequestID:  r.Context().Value(middlewares.RequestIDKey).(string),
		Endpoint:   r.RequestURI,
		Method:     r.Method,
		Operation:  "Update App Config",
		Status:     "Init",
		UserID:     strconv.FormatInt(userID, 10),
		Message:    "App config updation initiated.",
		ResourceID: strconv.Itoa(appID),
	}
	logger.Info(logdata)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logdata.Message = "Request body parsing failed."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, "Invalid request body. 'smtp_host', 'smtp_port', 'smtp_username', 'smtp_password', 'open_track', 'click_track' required.", logdata.RequestID)
		return
	}

	if err := services.ValidateAppConfigUpdate(req); err != nil {
		logdata.Message = "Request body validation failed."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, err.Error(), logdata.RequestID)
		return
	}

	exists, err := h.AppConfigRepo.AppConfigExistsByID(int64(appID), userID)
	if !exists {
		logdata.Message = "App config not found."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusNotFound
		logdata.Error = ""
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusNotFound, nil, "App config not found.", logdata.RequestID)
		return
	}
	if err != nil {
		logdata.Message = "App config existence check failed."
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	encryptedPassword, err := services.Encrypt(req.SMTPPassword)
	if err != nil {
		logdata.Message = "Password encryption failure"
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	appConfig := models.AppConfigUpdate{
		ID:           req.ID,
		SMTPHost:     req.SMTPHost,
		SMTPPort:     req.SMTPPort,
		SMTPUsername: req.SMTPUsername,
		SMTPPassword: encryptedPassword,
		OpenTrack:    req.OpenTrack,
		ClickTrack:   req.ClickTrack,
	}

	if err := h.AppConfigRepo.UpdateAppConfig(int64(appID), appConfig); err != nil {
		logdata.Message = "App config updation failed"
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	logdata.Message = "App config Updation successful."
	logdata.Status = "Success"
	logdata.ResponseCode = http.StatusOK
	logdata.Error = ""
	logger.Info(logdata)
	services.ResponseWithMessage(w, http.StatusOK, nil, "App config updated successfully.", logdata.RequestID)
}

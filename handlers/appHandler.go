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

	"github.com/google/uuid"
)

type AppHandler struct {
	AppRepo *repositories.AppRepository
}

func (h *AppHandler) CreateApp(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middlewares.UserIDKey).(int64)

	logdata := models.LogData{
		RequestID: r.Context().Value(middlewares.RequestIDKey).(string),
		Endpoint:  r.RequestURI,
		Method:    r.Method,
		Operation: "App Create",
		Status:    "Init",
		UserID:    strconv.FormatInt(userID, 10),
		Message:   "App Creation initiated.",
	}
	logger.Info(logdata)

	var req models.CreateApp

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logdata.Message = "Request body parsing failed."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, "Invalid request body. 'name', 'description', 'sender_name' and 'sender_email' required.", logdata.RequestID)
		return
	}

	if err := services.ValidateAppCreate(req); err != nil {
		logdata.Message = "Request body validation failed."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, err.Error(), logdata.RequestID)
		return
	}

	exists, err := h.AppRepo.AppExists(req.Name, userID)
	if err != nil {
		logdata.Message = "App name exist check failure"
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	if exists {
		logdata.Message = "App Name already exists."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = ""
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, "App name already exists.", logdata.RequestID)
		return
	}

	app := models.CreateApp{
		UserId:      userID,
		Name:        req.Name,
		Description: req.Description,
		SenderName:  req.SenderName,
		SenderEmail: req.SenderEmail,
		MailKey:     uuid.New(),
		Status:      "active",
	}

	if err := h.AppRepo.CreateApp(app); err != nil {
		logdata.Message = "App creation failure"
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	logdata.Message = "App creation successful."
	logdata.Status = "Success"
	logdata.ResponseCode = http.StatusOK
	logdata.Error = ""
	logger.Info(logdata)
	services.ResponseWithMessage(w, http.StatusCreated, nil, "App created successfully.", logdata.RequestID)
}

func (h *AppHandler) GetUserApps(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middlewares.UserIDKey).(int64)

	query := r.URL.Query()

	limit := 10
	if l := query.Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}

	offset := 0
	if o := query.Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}

	logdata := models.LogData{
		RequestID: r.Context().Value(middlewares.RequestIDKey).(string),
		Endpoint:  r.RequestURI,
		Method:    r.Method,
		Operation: "Fetch User Apps",
		Status:    "Init",
		UserID:    strconv.FormatInt(userID, 10),
		Message:   "User apps fetch initiated.",
	}
	logger.Info(logdata)

	apps, err := h.AppRepo.GetUserApps(userID, limit, offset)
	if err != nil {
		logdata.Message = "No users app fetched."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusNotFound
		logdata.Error = err.Error()
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusNotFound, nil, "No apps found.", logdata.RequestID)
		return
	}

	logdata.Message = "User apps fetch successful."
	logdata.Status = "Success"
	logdata.ResponseCode = http.StatusOK
	logdata.Error = ""
	logger.Info(logdata)
	services.ResponseWithData(w, http.StatusOK, nil, "Apps found successfully", apps, logdata.RequestID)

}

func (h *AppHandler) GetAppSingle(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middlewares.UserIDKey).(int64)
	idStr := r.PathValue("appID")
	appID, _ := strconv.Atoi(idStr)

	logdata := models.LogData{
		RequestID:  r.Context().Value(middlewares.RequestIDKey).(string),
		Endpoint:   r.RequestURI,
		Method:     r.Method,
		Operation:  "Fetch App Data",
		Status:     "Init",
		UserID:     strconv.FormatInt(userID, 10),
		Message:    "User apps fetch initiated.",
		ResourceID: strconv.Itoa(appID),
	}
	logger.Info(logdata)

	exists, err := h.AppRepo.AppExistsByID(int64(appID), userID)
	if !exists {
		logdata.Message = "App not found."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusNotFound
		logdata.Error = ""
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusNotFound, nil, "App not found.", logdata.RequestID)
		return
	}

	appData, err := h.AppRepo.GetUserAppSingle(int64(appID), userID)
	if err != nil {
		logdata.Message = "No app data found."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusNotFound
		logdata.Error = err.Error()
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusNotFound, nil, "No app data found.", logdata.RequestID)
		return
	}

	logdata.Message = "App data fetch successful."
	logdata.Status = "Success"
	logdata.ResponseCode = http.StatusOK
	logdata.Error = ""
	logger.Info(logdata)
	services.ResponseWithData(w, http.StatusOK, nil, "App found successfully", appData, logdata.RequestID)

}

func (h *AppHandler) GetAppKey(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middlewares.UserIDKey).(int64)
	idStr := r.PathValue("appID")
	appID, _ := strconv.Atoi(idStr)

	logdata := models.LogData{
		RequestID:  r.Context().Value(middlewares.RequestIDKey).(string),
		Endpoint:   r.RequestURI,
		Method:     r.Method,
		Operation:  "Fetch App Key",
		Status:     "Init",
		UserID:     strconv.FormatInt(userID, 10),
		Message:    "App Key fetch initiated.",
		ResourceID: strconv.Itoa(appID),
	}
	logger.Info(logdata)

	exists, err := h.AppRepo.AppExistsByID(int64(appID), userID)
	if !exists {
		logdata.Message = "App not found."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusNotFound
		logdata.Error = ""
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusNotFound, nil, "App not found.", logdata.RequestID)
		return
	}

	appData, err := h.AppRepo.GetUserAppKey(int64(appID), userID)
	if err != nil {
		logdata.Message = "App key fetch failed"
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	if appData == nil {
		logdata.Message = "No app key found."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusNotFound
		logdata.Error = ""
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusNotFound, nil, "No apps found.", logdata.RequestID)
		return
	}

	logdata.Message = "App Key fetch successful."
	logdata.Status = "Success"
	logdata.ResponseCode = http.StatusOK
	logdata.Error = ""
	logger.Info(logdata)

	services.ResponseWithData(w, http.StatusOK, nil, "Key found successfully", appData, logdata.RequestID)

}

func (h *AppHandler) UpdateApp(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middlewares.UserIDKey).(int64)
	idStr := r.PathValue("appID")
	appID, _ := strconv.Atoi(idStr)

	var req models.UpdateApp

	logdata := models.LogData{
		RequestID:  r.Context().Value(middlewares.RequestIDKey).(string),
		Endpoint:   r.RequestURI,
		Method:     r.Method,
		Operation:  "Update App",
		Status:     "Init",
		UserID:     strconv.FormatInt(userID, 10),
		Message:    "App updation initiated.",
		ResourceID: strconv.Itoa(appID),
	}
	logger.Info(logdata)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logdata.Message = "Request body parsing failed."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, "Invalid request body. 'name', 'description', 'sender_name' and 'sender_email' required.", logdata.RequestID)
		return
	}

	if err := services.ValidateAppUpdate(req); err != nil {
		logdata.Message = "Request body validation failed."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, err.Error(), logdata.RequestID)
		return
	}

	exists, err := h.AppRepo.AppExistsByID(int64(appID), userID)
	if !exists {
		logdata.Message = "App not found."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusNotFound
		logdata.Error = ""
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusNotFound, nil, "App not found.", logdata.RequestID)
		return
	}
	if err != nil {
		logdata.Message = "App existence check failed."
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	app := models.UpdateApp{
		UserId:      userID,
		Name:        req.Name,
		Description: req.Description,
		SenderName:  req.SenderName,
		SenderEmail: req.SenderEmail,
	}

	if err := h.AppRepo.UpdateApp(int64(appID), app); err != nil {
		logdata.Message = "App updation failed"
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	logdata.Message = "App Updation successful."
	logdata.Status = "Success"
	logdata.ResponseCode = http.StatusOK
	logdata.Error = ""
	logger.Info(logdata)
	services.ResponseWithMessage(w, http.StatusOK, nil, "App updated successfully.", logdata.RequestID)
}

func (h *AppHandler) UpdateKey(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middlewares.UserIDKey).(int64)
	idStr := r.PathValue("appID")
	appID, _ := strconv.Atoi(idStr)

	logdata := models.LogData{
		RequestID:  r.Context().Value(middlewares.RequestIDKey).(string),
		Endpoint:   r.RequestURI,
		Method:     r.Method,
		Operation:  "Update Key",
		Status:     "Init",
		UserID:     strconv.FormatInt(userID, 10),
		Message:    "App key update initiated.",
		ResourceID: strconv.Itoa(appID),
	}
	logger.Info(logdata)

	exists, err := h.AppRepo.AppExistsByID(int64(appID), userID)
	if !exists {
		logdata.Message = "App not found."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusNotFound
		logdata.Error = ""
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusNotFound, nil, "App not found.", logdata.RequestID)
		return
	}
	if err != nil {
		logdata.Message = "App existence check failed."
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	mailKey := uuid.New()

	if err := h.AppRepo.RefreshMailKey(int64(appID), userID, mailKey); err != nil {
		logdata.Message = "App key update failed"
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	logdata.Message = "App key update successful."
	logdata.Status = "Success"
	logdata.ResponseCode = http.StatusOK
	logdata.Error = ""
	logger.Info(logdata)
	services.ResponseWithMessage(w, http.StatusOK, nil, "App key updated successfully.", logdata.RequestID)
}

func (h *AppHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middlewares.UserIDKey).(int64)
	idStr := r.PathValue("appID")
	appID, _ := strconv.Atoi(idStr)
	appStatus := r.PathValue("status")

	logdata := models.LogData{
		RequestID:  r.Context().Value(middlewares.RequestIDKey).(string),
		Endpoint:   r.RequestURI,
		Method:     r.Method,
		Operation:  "Update Status",
		Status:     "Init",
		UserID:     strconv.FormatInt(userID, 10),
		Message:    "App status update initiated.",
		ResourceID: strconv.Itoa(appID),
	}
	logger.Info(logdata)

	exists, err := h.AppRepo.AppExistsByID(int64(appID), userID)
	if !exists {
		logdata.Message = "App not found."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusNotFound
		logdata.Error = ""
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusNotFound, nil, "App not found.", logdata.RequestID)
		return
	}
	if err != nil {
		logdata.Message = "App existence check failed."
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	mailKey := uuid.New()

	if appStatus == "active" {
		mailKey = uuid.New()
	} else if appStatus == "inactive" {
		mailKey = uuid.Nil
	} else {
		logdata.Message = "Invalid status value."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = "Status must be either 'active' or 'inactive'."
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, "Invalid status value. Status must be either 'active' or 'inactive'.", logdata.RequestID)
		return
	}

	if err := h.AppRepo.UpdateAppStatus(int64(appID), userID, appStatus, mailKey); err != nil {
		logdata.Message = "App status update failed"
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	logdata.Message = "App status update successful."
	logdata.Status = "Success"
	logdata.ResponseCode = http.StatusOK
	logdata.Error = ""
	logger.Info(logdata)
	services.ResponseWithMessage(w, http.StatusOK, nil, "App status updated successfully.", logdata.RequestID)
}

func (h *AppHandler) DeleteApp(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middlewares.UserIDKey).(int64)
	idStr := r.PathValue("appID")
	appID, _ := strconv.Atoi(idStr)

	logdata := models.LogData{
		RequestID:  r.Context().Value(middlewares.RequestIDKey).(string),
		Endpoint:   r.RequestURI,
		Method:     r.Method,
		Operation:  "Delete App",
		Status:     "Init",
		UserID:     strconv.FormatInt(userID, 10),
		Message:    "App Deletion initiated.",
		ResourceID: strconv.Itoa(appID),
	}
	logger.Info(logdata)

	exists, err := h.AppRepo.AppExistsByID(int64(appID), userID)
	if !exists {
		logdata.Message = "App not found."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusNotFound
		logdata.Error = ""
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusNotFound, nil, "App not found.", logdata.RequestID)
		return
	}
	if err != nil {
		logdata.Message = "App existence check failed."
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	if err := h.AppRepo.DeleteApp(int64(appID), userID); err != nil {
		logdata.Message = "App deletion failed"
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	logdata.Message = "App deletion successful."
	logdata.Status = "Success"
	logdata.ResponseCode = http.StatusOK
	logdata.Error = ""
	logger.Info(logdata)

	services.ResponseWithMessage(w, http.StatusOK, nil, "App deleted successfully.", logdata.RequestID)
}

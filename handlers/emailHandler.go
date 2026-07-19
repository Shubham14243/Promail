package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"promail/logger"
	"promail/middlewares"
	"promail/models"
	"promail/repositories"
	"promail/services"
)

type EmailHandler struct {
	EmailRepo     *repositories.EmailRepository
	AppConfigRepo *repositories.AppConfigRepository
	AppRepo       *repositories.AppRepository
	TempRepo      *repositories.TemplateRepository
}

func (h *EmailHandler) SendEmailTest(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middlewares.UserIDKey).(int64)

	logdata := models.LogData{
		RequestID: r.Context().Value(middlewares.RequestIDKey).(string),
		Endpoint:  r.RequestURI,
		Method:    r.Method,
		Operation: "Email Send Test Request.",
		Status:    "Init",
		UserID:    strconv.FormatInt(userID, 10),
		Message:   "Email sending initiated.",
	}
	logger.Info(logdata)

	var req models.EmailSendTest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logdata.Message = "Request body parsing failed."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, "Invalid request body. 'app_id', 'mail_key', 'to', 'subject', 'body' required.", logdata.RequestID)
		return
	}

	if err := services.ValidateEmailTestData(req); err != nil {
		logdata.Message = "Request body validation failed."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, err.Error(), logdata.RequestID)
		return
	}

	exists, err := h.AppConfigRepo.AppConfigExistsByAppID(int64(req.AppID), int64(userID))
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
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	appConf.SMTPPassword = decrypted_password

	app_mailkey, err := h.AppRepo.GetUserAppKey(int64(req.AppID), userID)
	if err != nil {
		logdata.Message = "User app key not found."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusNotFound
		logdata.Error = err.Error()
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusNotFound, nil, "User app key not found.", logdata.RequestID)
		return
	}

	if app_mailkey.MailKey.String() != req.MailKey {
		logdata.Message = "Invalid mail_key provided."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = ""
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, "Invalid or expired mail_key provided.", logdata.RequestID)
		return
	}

	if err := services.SendEmail(appConf, req.To, req.Subject, req.Body, "html"); err != nil {
		logdata.Message = "Email sending test failure."
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	logdata.Message = "Email sending test successful."
	logdata.Status = "Success"
	logdata.ResponseCode = http.StatusOK
	logdata.Error = ""
	logger.Info(logdata)
	services.ResponseWithMessage(w, http.StatusOK, nil, "Email test sent successfully.", logdata.RequestID)

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
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, "Invalid request body. 'app_id', 'template_slug', 'mail_key', 'to', 'variables' required.", logdata.RequestID)
		return
	}

	if err := services.ValidateEmailData(req); err != nil {
		logdata.Message = "Request body validation failed."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, err.Error(), logdata.RequestID)
		return
	}
	logdata.ResourceID = strconv.Itoa(int(req.AppID))

	exists, err := h.AppConfigRepo.AppConfigExistsByAppID(int64(req.AppID), int64(userID))
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
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	appConf.SMTPPassword = decrypted_password

	app_mailkey, err := h.AppRepo.GetUserAppKey(int64(req.AppID), userID)
	if err != nil {
		logdata.Message = "User app key not found."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = err.Error()
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, "User app key not found.", logdata.RequestID)
		return
	}

	if app_mailkey.MailKey.String() != req.MailKey {
		logdata.Message = "Invalid mail_key provided."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusUnauthorized
		logdata.Error = ""
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusUnauthorized, nil, "Invalid or expired mail_key provided.", logdata.RequestID)
		return
	}

	template, err := h.TempRepo.GetAppTemplateBySlug(req.TemplateSlug, int64(userID))
	if err != nil {
		logdata.Message = "Template fetch failure."
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	var varData []byte
	varData, err = json.Marshal(req.Variables)
	if err != nil {
		logdata.Message = "Variable data string conversion failed."
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	logUUID := uuid.New()

	email_body := services.PrepareEmailBody(template.Content, req.Variables)

	var trackings []models.ClickTracking

	if appConf.ClickTrack == "active" {
		email_body, trackings = services.AddClickTracking(email_body)
	}

	if appConf.OpenTrack == "active" {
		email_body = services.AddOpenTracking(email_body, logUUID.String(), template.Type)
	}

	emailLog := models.EmailLogDataCreate{
		UUID:         logUUID,
		UserID:       int64(userID),
		AppID:        int64(req.AppID),
		TemplateID:   template.ID,
		ToEmail:      req.To,
		Subject:      template.Subject,
		VariableData: string(varData),
		Body:         email_body,
	}

	var logRes models.LogResponse
	logRes, err = h.EmailRepo.AddEmailLogAndQueue(emailLog)
	if err != nil {
		logdata.Message = "Email logging and queueing failed."
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Failed to log and queue Email.", logdata.RequestID)
		return
	}

	if appConf.OpenTrack == "active" {
		if err := h.EmailRepo.AddOpenTracking(logRes); err != nil {
			logdata.Message = "Open tracking data insertion failed."
			logdata.Status = "Failure"
			logdata.ResponseCode = http.StatusBadRequest
			logdata.Error = err.Error()
			logger.Error(logdata)
			services.ResponseWithMessage(w, http.StatusBadRequest, nil, "Open tracking data insertion failed.", logdata.RequestID)
			return
		}
	}

	if appConf.ClickTrack == "active" {
		err = h.EmailRepo.AddClickTracking(logRes.LogID, trackings)
		if err != nil {
			logdata.Message = "Click tracking creation failed."
			logdata.Status = "Error"
			logdata.ResponseCode = http.StatusInternalServerError
			logdata.Error = err.Error()
			logger.Error(logdata)
			services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
			return
		}
	}

	logdata.Message = "Email accepted successful."
	logdata.Status = "Success"
	logdata.ResponseCode = http.StatusAccepted
	logdata.Error = ""
	logger.Info(logdata)
	services.ResponseWithData(w, http.StatusAccepted, nil, "Email accepted successfully.", map[string]uuid.UUID{"acknowledgement_id": logRes.AckID}, logdata.RequestID)

}

func (h *EmailHandler) OpenTrack(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.PathValue("token")

	logdata := models.LogData{
		RequestID:  r.Context().Value(middlewares.RequestIDKey).(string),
		Endpoint:   r.RequestURI,
		Method:     r.Method,
		Operation:  "Open Tracking",
		Status:     "Init",
		UserID:     "",
		Message:    "Email open tracking initiated.",
		ResourceID: tokenStr,
	}
	logger.Info(logdata)

	trackData, err := h.EmailRepo.GetAnalyticsWithUUID(tokenStr)
	if err != nil {
		logdata.Message = "LogData fetch failure."
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	if trackData == nil {
		logdata.Message = "Open tracking record not found."
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusNotFound
		logdata.Error = ""
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusNotFound, nil, "Open tracking record not found.", logdata.RequestID)
		return
	}

	if trackData.OpenedAt == nil {
		if err := h.EmailRepo.UpdateOpenTracking(trackData.EmailLogID); err != nil {
			logdata.Message = "Open tracking data updation failed."
			logdata.Status = "Failure"
			logdata.ResponseCode = http.StatusBadRequest
			logdata.Error = err.Error()
			logger.Error(logdata)
			services.ResponseWithMessage(w, http.StatusBadRequest, nil, "Open tracking data updation failed.", logdata.RequestID)
			return
		}
	}

	logdata.Message = "Open Track Successful."
	logdata.Status = "Success"
	logdata.ResponseCode = http.StatusOK
	logdata.Error = ""
	logger.Info(logdata)
	services.ResponseWithMessage(w, http.StatusOK, nil, "Open Track Successful.", logdata.RequestID)
}

func (h *EmailHandler) ClickTrack(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.PathValue("token")

	logdata := models.LogData{
		RequestID:  r.Context().Value(middlewares.RequestIDKey).(string),
		Endpoint:   r.RequestURI,
		Method:     r.Method,
		Operation:  "Open Tracking",
		Status:     "Init",
		UserID:     "",
		Message:    "Email open tracking initiated.",
		ResourceID: tokenStr,
	}
	logger.Info(logdata)

	trackData, err := h.EmailRepo.GetAnalyticsWithUUID(tokenStr)
	if err != nil {
		logdata.Message = "LogData fetch failure."
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	if trackData == nil {
		logdata.Message = "Click tracking record not found."
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusNotFound
		logdata.Error = ""
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusNotFound, nil, "Click tracking record not found.", logdata.RequestID)
		return
	}

	if trackData.ClickedAt == nil {
		if err := h.EmailRepo.UpdateClickTracking(trackData.EmailLogID, tokenStr); err != nil {
			logdata.Message = "Click tracking data updation failed."
			logdata.Status = "Failure"
			logdata.ResponseCode = http.StatusBadRequest
			logdata.Error = err.Error()
			logger.Error(logdata)
			services.ResponseWithMessage(w, http.StatusBadRequest, nil, "Click tracking data updation failed.", logdata.RequestID)
			return
		}
	}

	logdata.Message = "Click Track Successful."
	logdata.Status = "Success"
	logdata.ResponseCode = http.StatusOK
	logdata.Error = ""
	logger.Info(logdata)
	services.ResponseWithMessage(w, http.StatusOK, nil, "Click Track Successful.", logdata.RequestID)
}

func (h *EmailHandler) EmailLogUUID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.UserIDKey).(int64)
	uuidStr := r.PathValue("uuid")

	logdata := models.LogData{
		RequestID:  r.Context().Value(middlewares.RequestIDKey).(string),
		Endpoint:   r.RequestURI,
		Method:     r.Method,
		Operation:  "Email Log Fetch",
		Status:     "Init",
		UserID:     strconv.FormatInt(userID, 10),
		Message:    "Email log fetching initiated.",
		ResourceID: uuidStr,
	}
	logger.Info(logdata)

	emailLog, err := h.EmailRepo.GetLogWithUUID(uuidStr)
	if err != nil {
		logdata.Message = "Email log fetch failure."
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	if emailLog == nil {
		logdata.Message = "Email Log record not found."
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusNotFound
		logdata.Error = ""
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusNotFound, nil, "Email Log record not found.", logdata.RequestID)
		return
	}

	logdata.Message = "Email log fetch Successful."
	logdata.Status = "Success"
	logdata.ResponseCode = http.StatusOK
	logdata.Error = ""
	logger.Info(logdata)
	services.ResponseWithData(w, http.StatusOK, nil, "Email log fetch Successful.", emailLog, logdata.RequestID)
}

func (h *EmailHandler) EmailLogAll(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middlewares.UserIDKey).(int64)

	logdata := models.LogData{
		RequestID:  r.Context().Value(middlewares.RequestIDKey).(string),
		Endpoint:   r.RequestURI,
		Method:     r.Method,
		Operation:  "Email Log Fetch",
		Status:     "Init",
		UserID:     strconv.FormatInt(userID, 10),
		Message:    "Email log fetching initiated.",
		ResourceID: "",
	}
	logger.Info(logdata)

	query := r.URL.Query()

	var filter models.EmailLogFilter

	if appID := query.Get("app_id"); appID != "" {
		id, _ := strconv.ParseInt(appID, 10, 64)
		filter.AppID = &id
	}

	if templateID := query.Get("template_id"); templateID != "" {
		id, _ := strconv.ParseInt(templateID, 10, 64)
		filter.TemplateID = &id
	}

	filter.ToEmail = query.Get("to_email")
	filter.StartDateTime = query.Get("startDateTime")
	filter.EndDateTime = query.Get("endDateTime")

	emailLogs, err := h.EmailRepo.GetEmailLogs(userID, filter)
	if err != nil {
		logdata.Message = "Email log fetch failure."
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	if emailLogs == nil {
		logdata.Message = "Email Log record not found."
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusNotFound
		logdata.Error = ""
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusNotFound, nil, "Email Log record not found.", logdata.RequestID)
		return
	}

	logdata.Message = "Email log fetch Successful."
	logdata.Status = "Success"
	logdata.ResponseCode = http.StatusOK
	logdata.Error = ""
	logger.Info(logdata)
	services.ResponseWithData(w, http.StatusOK, nil, "Email log fetch Successful.", emailLogs, logdata.RequestID)
}

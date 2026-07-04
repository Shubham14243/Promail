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

type TemplateHandler struct {
	TempRepo *repositories.TemplateRepository
}

func (h *TemplateHandler) CreateTemplate(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middlewares.UserIDKey).(int64)

	logdata := models.LogData{
		RequestID: r.Context().Value(middlewares.RequestIDKey).(string),
		Endpoint:  r.RequestURI,
		Method:    r.Method,
		Operation: "Template Create",
		Status:    "Init",
		UserID:    strconv.FormatInt(userID, 10),
		Message:   "Template Creation initiated.",
	}
	logger.Info(logdata)

	var req models.TemplateCreate

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logdata.Message = "Request body parsing failed."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, "Invalid request body. 'app_id', 'name', 'slug', 'subject', 'type' and 'content' required.", logdata.RequestID)
		return
	}

	if err := services.ValidateTemplateCreate(req); err != nil {
		logdata.Message = "Request body validation failed."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, err.Error(), logdata.RequestID)
		return
	}

	exists, err := h.TempRepo.TemplateExistsBySlug(req.Slug, userID)
	if err != nil {
		logdata.Message = "Template slug exist check failure"
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	if exists {
		logdata.Message = "Template slug already exists."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = ""
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, "Template slug already exists.", logdata.RequestID)
		return
	}

	template := models.TemplateCreate{
		AppID:   req.AppID,
		Name:    req.Name,
		Slug:    req.Slug,
		Subject: req.Subject,
		Type:    req.Type,
		Content: req.Content,
		Status:  "active",
	}

	if err := h.TempRepo.CreateTemplate(template); err != nil {
		logdata.Message = "Template creation failure"
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	logdata.Message = "Template creation successful."
	logdata.Status = "Success"
	logdata.ResponseCode = http.StatusOK
	logdata.Error = ""
	logger.Info(logdata)
	services.ResponseWithMessage(w, http.StatusCreated, nil, "Template created successfully.", logdata.RequestID)
}

func (h *TemplateHandler) GetAppTemplates(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middlewares.UserIDKey).(int64)
	idStr := r.PathValue("appID")
	appID, _ := strconv.Atoi(idStr)

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
		Operation: "Fetch App Templates",
		Status:    "Init",
		UserID:    strconv.FormatInt(userID, 10),
		Message:   "User templates fetch initiated.",
	}
	logger.Info(logdata)

	templates, err := h.TempRepo.GetAppTemplates(int64(appID), userID, limit, offset)
	if err != nil {
		logdata.Message = "No app template fetched."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusNotFound
		logdata.Error = err.Error()
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusNotFound, nil, "No templates found.", logdata.RequestID)
		return
	}

	logdata.Message = "User templates fetch successful."
	logdata.Status = "Success"
	logdata.ResponseCode = http.StatusOK
	logdata.Error = ""
	logger.Info(logdata)
	services.ResponseWithData(w, http.StatusOK, nil, "Templates found successfully", templates, logdata.RequestID)

}

func (h *TemplateHandler) GetTemplateData(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middlewares.UserIDKey).(int64)
	idStr := r.PathValue("templateID")
	templateID, _ := strconv.Atoi(idStr)

	logdata := models.LogData{
		RequestID:  r.Context().Value(middlewares.RequestIDKey).(string),
		Endpoint:   r.RequestURI,
		Method:     r.Method,
		Operation:  "Fetch App Data",
		Status:     "Init",
		UserID:     strconv.FormatInt(userID, 10),
		Message:    "User templates fetch initiated.",
		ResourceID: strconv.Itoa(templateID),
	}
	logger.Info(logdata)

	exists, err := h.TempRepo.TemplateExistsByID(int64(templateID), userID)
	if !exists {
		logdata.Message = "Template not found."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusNotFound
		logdata.Error = ""
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusNotFound, nil, "Template not found.", logdata.RequestID)
		return
	}

	templateData, err := h.TempRepo.GetAppTemplateSingle(int64(templateID), userID)
	if err != nil {
		logdata.Message = "No template data found."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusNotFound
		logdata.Error = err.Error()
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusNotFound, nil, "No template data found.", logdata.RequestID)
		return
	}

	logdata.Message = "Template data fetch successful."
	logdata.Status = "Success"
	logdata.ResponseCode = http.StatusOK
	logdata.Error = ""
	logger.Info(logdata)
	services.ResponseWithData(w, http.StatusOK, nil, "Template found successfully", templateData, logdata.RequestID)

}

func (h *TemplateHandler) UpdateTemplate(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middlewares.UserIDKey).(int64)
	idStr := r.PathValue("templateID")
	templateID, _ := strconv.Atoi(idStr)

	var req models.TemplateUpdate

	logdata := models.LogData{
		RequestID:  r.Context().Value(middlewares.RequestIDKey).(string),
		Endpoint:   r.RequestURI,
		Method:     r.Method,
		Operation:  "Update Template",
		Status:     "Init",
		UserID:     strconv.FormatInt(userID, 10),
		Message:    "Template updation initiated.",
		ResourceID: strconv.Itoa(templateID),
	}
	logger.Info(logdata)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logdata.Message = "Request body parsing failed."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, "Invalid request body. 'name', 'slug', 'subject' and 'status' required.", logdata.RequestID)
		return
	}

	if err := services.ValidateTemplateUpdate(req); err != nil {
		logdata.Message = "Request body validation failed."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, err.Error(), logdata.RequestID)
		return
	}

	exists, err := h.TempRepo.TemplateExistsByID(int64(templateID), userID)
	if !exists {
		logdata.Message = "Template not found."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusNotFound
		logdata.Error = ""
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusNotFound, nil, "Template not found.", logdata.RequestID)
		return
	}
	if err != nil {
		logdata.Message = "Template existence check failed."
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	template := models.TemplateUpdate{
		Name:    req.Name,
		Slug:    req.Slug,
		Subject: req.Subject,
		Status:  req.Status,
	}

	if err := h.TempRepo.UpdateTemplate(int64(templateID), template); err != nil {
		logdata.Message = "Template updation failed"
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	logdata.Message = "Template Updation successful."
	logdata.Status = "Success"
	logdata.ResponseCode = http.StatusOK
	logdata.Error = ""
	logger.Info(logdata)
	services.ResponseWithMessage(w, http.StatusOK, nil, "Template updated successfully.", logdata.RequestID)
}

func (h *TemplateHandler) UpdateContent(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middlewares.UserIDKey).(int64)
	idStr := r.PathValue("templateID")
	templateID, _ := strconv.Atoi(idStr)

	var req models.TemplateContent

	logdata := models.LogData{
		RequestID:  r.Context().Value(middlewares.RequestIDKey).(string),
		Endpoint:   r.RequestURI,
		Method:     r.Method,
		Operation:  "Update Content",
		Status:     "Init",
		UserID:     strconv.FormatInt(userID, 10),
		Message:    "Template content update initiated.",
		ResourceID: strconv.Itoa(templateID),
	}
	logger.Info(logdata)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logdata.Message = "Request body parsing failed."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, "Invalid request body. 'type' and 'content' required.", logdata.RequestID)
		return
	}

	if err := services.ValidateTemplateContent(req); err != nil {
		logdata.Message = "Request body validation failed."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, err.Error(), logdata.RequestID)
		return
	}

	exists, err := h.TempRepo.TemplateExistsByID(int64(templateID), userID)
	if !exists {
		logdata.Message = "Template not found."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusNotFound
		logdata.Error = ""
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusNotFound, nil, "Template not found.", logdata.RequestID)
		return
	}
	if err != nil {
		logdata.Message = "Template existence check failed."
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	templateContent := models.TemplateContent{
		Type:    req.Type,
		Content: req.Content,
	}

	if err := h.TempRepo.UpdateTemplateContent(int64(templateID), templateContent); err != nil {
		logdata.Message = "Template content update failed"
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	logdata.Message = "Template content update successful."
	logdata.Status = "Success"
	logdata.ResponseCode = http.StatusOK
	logdata.Error = ""
	logger.Info(logdata)
	services.ResponseWithMessage(w, http.StatusOK, nil, "Template content updated successfully.", logdata.RequestID)
}

func (h *TemplateHandler) DeleteTemplate(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middlewares.UserIDKey).(int64)
	idStr := r.PathValue("templateID")
	templateID, _ := strconv.Atoi(idStr)

	logdata := models.LogData{
		RequestID:  r.Context().Value(middlewares.RequestIDKey).(string),
		Endpoint:   r.RequestURI,
		Method:     r.Method,
		Operation:  "Delete Template",
		Status:     "Init",
		UserID:     strconv.FormatInt(userID, 10),
		Message:    "Template Deletion initiated.",
		ResourceID: strconv.Itoa(templateID),
	}
	logger.Info(logdata)

	exists, err := h.TempRepo.TemplateExistsByID(int64(templateID), userID)
	if !exists {
		logdata.Message = "Template not found."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusNotFound
		logdata.Error = ""
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusNotFound, nil, "Template not found.", logdata.RequestID)
		return
	}
	if err != nil {
		logdata.Message = "Template existence check failed."
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	if err := h.TempRepo.DeleteTemplate(int64(templateID), userID); err != nil {
		logdata.Message = "Template deletion failed"
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Info(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	logdata.Message = "Template deletion successful."
	logdata.Status = "Success"
	logdata.ResponseCode = http.StatusOK
	logdata.Error = ""
	logger.Info(logdata)
	services.ResponseWithMessage(w, http.StatusOK, nil, "Template deleted successfully.", logdata.RequestID)
}

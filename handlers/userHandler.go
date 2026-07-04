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

type UserHandler struct {
	Repo *repositories.UserRepository
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {

	logdata := models.LogData{
		RequestID: r.Context().Value(middlewares.RequestIDKey).(string),
		Endpoint:  r.RequestURI,
		Method:    r.Method,
		Operation: "User GetUsers",
		Status:    "Init",
		UserID:    strconv.FormatInt(r.Context().Value(middlewares.UserIDKey).(int64), 10),
		Message:   "Users Fetch initiated.",
	}
	logger.Info(logdata)

	users, err := h.Repo.GetAllUsers()
	if err != nil {
		logdata.Message = "Users data fetch failed"
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, err.Error(), logdata.RequestID)
		return
	}

	if users == nil {
		logdata.Message = "No users data fetched."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusNotFound
		logdata.Error = ""
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusNotFound, nil, "No users found.", logdata.RequestID)
		return
	}

	logdata.Message = "Users data successful."
	logdata.Status = "Success"
	logdata.ResponseCode = http.StatusOK
	logdata.Error = ""
	logger.Info(logdata)
	services.ResponseWithData(w, http.StatusOK, nil, "Users found successfully.", users, logdata.RequestID)

}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middlewares.UserIDKey).(int64)

	logdata := models.LogData{
		RequestID: r.Context().Value(middlewares.RequestIDKey).(string),
		Endpoint:  r.RequestURI,
		Method:    r.Method,
		Operation: "User GetUser Single",
		Status:    "Init",
		UserID:    strconv.FormatInt(userID, 10),
		Message:   "User single Fetch initiated.",
	}
	logger.Info(logdata)

	user, err := h.Repo.GetUserByID(userID)
	if err != nil {
		logdata.Message = "No single user data fetched."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusNotFound
		logdata.Error = ""
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusNotFound, nil, "User not found.", logdata.RequestID)
		return
	}

	logdata.Message = "User single data successful."
	logdata.Status = "Success"
	logdata.ResponseCode = http.StatusOK
	logdata.Error = ""
	logger.Info(logdata)
	services.ResponseWithData(w, http.StatusOK, nil, "User found successfully.", user, logdata.RequestID)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middlewares.UserIDKey).(int64)

	logdata := models.LogData{
		RequestID:  r.Context().Value(middlewares.RequestIDKey).(string),
		Endpoint:   r.RequestURI,
		Method:     r.Method,
		Operation:  "Update User",
		Status:     "Init",
		UserID:     strconv.FormatInt(userID, 10),
		Message:    "User updation initiated.",
		ResourceID: strconv.FormatInt(userID, 10),
	}
	logger.Info(logdata)

	var req models.UserUpdateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logdata.Message = "Request body parsing failed."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, "Invalid request body. 'name' and 'email' required.", logdata.RequestID)
		return
	}

	if err := services.ValidateUserUpdate(req); err != nil {
		logdata.Message = "Request body validation failed."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, err.Error(), logdata.RequestID)
		return
	}

	if err := h.Repo.UpdateUser(userID, req); err != nil {
		logdata.Message = "User updation failed"
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	logdata.Message = "User Updation successful."
	logdata.Status = "Success"
	logdata.ResponseCode = http.StatusOK
	logdata.Error = ""
	logger.Info(logdata)
	services.ResponseWithMessage(w, http.StatusOK, nil, "User updated successfully.", logdata.RequestID)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middlewares.UserIDKey).(int64)

	logdata := models.LogData{
		RequestID:  r.Context().Value(middlewares.RequestIDKey).(string),
		Endpoint:   r.RequestURI,
		Method:     r.Method,
		Operation:  "Delete User",
		Status:     "Init",
		UserID:     strconv.FormatInt(userID, 10),
		Message:    "User Deletion initiated.",
		ResourceID: strconv.FormatInt(userID, 10),
	}
	logger.Info(logdata)

	if err := h.Repo.DeleteUser(userID); err != nil {
		logdata.Message = "User deletion failed"
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	logdata.Message = "User deletion successful."
	logdata.Status = "Success"
	logdata.ResponseCode = http.StatusOK
	logdata.Error = ""
	logger.Info(logdata)
	services.ResponseWithMessage(w, http.StatusOK, nil, "User deleted successfully.", logdata.RequestID)
}

package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"promail/logger"
	"promail/middlewares"
	"promail/models"
	"promail/repositories"
	"promail/services"
)

type AuthHandler struct {
	UserRepo         *repositories.UserRepository
	RefreshTokenRepo *repositories.RefreshTokenRepository
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {

	logdata := models.LogData{
		RequestID: r.Context().Value(middlewares.RequestIDKey).(string),
		Endpoint:  r.RequestURI,
		Method:    r.Method,
		Operation: "User Signup",
		Status:    "Init",
		UserID:    "",
		Message:   "User signup initiated.",
	}
	logger.Info(logdata)

	var req models.UserCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logdata.Message = "Request body parsing failed."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, "Invalid request body. 'name', 'email' and 'password' required.", logdata.RequestID)
		return
	}

	if err := services.ValidateUserCreate(req); err != nil {
		logdata.Message = "Request body validation failed."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, err.Error(), logdata.RequestID)
		return
	}

	exists, err := h.UserRepo.UserExists(req.Email)
	if err != nil {
		logdata.Message = "User exist check failure"
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	if exists {
		logdata.Message = "User already exists."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = ""
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, "User already exists.", logdata.RequestID)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logdata.Message = "Password hashing failure"
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	user := models.User{
		UUID:         uuid.New(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := h.UserRepo.CreateUser(user); err != nil {
		logdata.Message = "User creation failure"
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	logdata.Message = "User creation successful."
	logdata.Status = "Success"
	logdata.ResponseCode = http.StatusOK
	logdata.Error = ""
	logger.Info(logdata)
	services.ResponseWithMessage(w, http.StatusCreated, nil, "User created successfully.", logdata.RequestID)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	logdata := models.LogData{
		RequestID: r.Context().Value(middlewares.RequestIDKey).(string),
		Endpoint:  r.RequestURI,
		Method:    r.Method,
		Operation: "User Login",
		Status:    "Init",
		UserID:    "",
		Message:   "User login initiated.",
	}
	logger.Info(logdata)

	var req models.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logdata.Message = "Request body parsing failed."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, "Invalid request body. 'name', 'email' and 'password' required.", logdata.RequestID)
		return
	}

	if err := services.ValidateLogin(req); err != nil {
		logdata.Message = "Request body validation failed."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, err.Error(), logdata.RequestID)
		return
	}

	user, err := h.UserRepo.GetUserByEmail(req.Email)

	if err != nil {
		logdata.Message = "No User found with - " + req.Email
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusNotFound, nil, "Invalid email or passowrd.", logdata.RequestID)
		return
	}

	logdata.UserID = strconv.Itoa(int(user.ID))

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		logdata.Message = "Invalid Password. User - " + req.Email
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusNotFound, nil, "Invalid email or passowrd.", logdata.RequestID)
		return
	}

	authToken, err := services.GenerateAccessToken(user.ID, req.Email)
	if err != nil {
		logdata.Message = "Login auth token generation failure"
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	refreshToken := services.GenerateRefreshToken()

	userData, err := h.UserRepo.GetUserByID(int64(user.ID))
	if err != nil {
		logdata.Message = "User data fetch failure"
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	validityDays, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_VALIDITY_DAYS"))
	if err != nil {
		validityDays = 3
	}

	refToken := models.RefreshTokenCreate{
		UserID:    userData.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Duration(validityDays) * 24 * time.Hour),
	}

	if err := h.RefreshTokenRepo.CreateRefreshToken(refToken); err != nil {
		logdata.Message = "Refresh token creation failure"
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	loginRes := models.LoginResponse{
		ID:           userData.ID,
		UUID:         userData.UUID,
		Name:         userData.Name,
		Email:        userData.Email,
		AuthToken:    authToken,
		RefreshToken: refreshToken,
		CreatedAt:    userData.CreatedAt,
		UpdatedAt:    userData.UpdatedAt,
	}

	logdata.Message = "User login successful."
	logdata.Status = "Success"
	logdata.ResponseCode = http.StatusOK
	logdata.Error = ""
	logger.Info(logdata)
	services.ResponseWithData(w, http.StatusOK, nil, "User loggedin successfully.", loginRes, logdata.RequestID)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {

	logdata := models.LogData{
		RequestID: r.Context().Value(middlewares.RequestIDKey).(string),
		Endpoint:  r.RequestURI,
		Method:    r.Method,
		Operation: "User Refresh Token",
		Status:    "Init",
		UserID:    "",
		Message:   "User refresh token initiated.",
	}
	logger.Info(logdata)

	var req models.RefreshTokenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logdata.Message = "Request body parsing failed."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, "Invalid request body. 'refresh_token' required.", logdata.RequestID)
		return
	}

	if req.RefreshToken == "" {
		logdata.Message = "Request body validation failed."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusBadRequest
		logdata.Error = "Refresh token empty."
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusBadRequest, nil, "refresh_token is required.", logdata.RequestID)
		return
	}

	rt, err := h.RefreshTokenRepo.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		logdata.Message = "Refresh token expired."
		logdata.Status = "Failure"
		logdata.ResponseCode = http.StatusUnauthorized
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusUnauthorized, nil, "invalid or expired refresh token", logdata.RequestID)
		return
	}

	userData, err := h.UserRepo.GetUserByID(int64(rt.UserID))
	logdata.UserID = strconv.Itoa(int(userData.ID))
	if err != nil {
		logdata.Message = "User data fetch failure"
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	accessToken, err := services.GenerateAccessToken(rt.UserID, userData.Email)
	if err != nil {
		logdata.Message = "Refresh auth token generation failure"
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "failed to generate access token", logdata.RequestID)
		return
	}

	resp := models.RefreshTokenResponse{
		AuthToken: accessToken,
	}

	logdata.Message = "Token refresh successful."
	logdata.Status = "Success"
	logdata.ResponseCode = http.StatusOK
	logdata.Error = ""
	logger.Info(logdata)
	services.ResponseWithData(w, http.StatusOK, nil, "token refreshed successfully", resp, logdata.RequestID)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middlewares.UserIDKey).(int64)

	logdata := models.LogData{
		RequestID: r.Context().Value(middlewares.RequestIDKey).(string),
		Endpoint:  r.RequestURI,
		Method:    r.Method,
		Operation: "User Logout",
		Status:    "Init",
		UserID:    strconv.FormatInt(userID, 10),
		Message:   "Logout initiated.",
	}
	logger.Info(logdata)

	if err := h.RefreshTokenRepo.DeleteRefreshToken(userID); err != nil {
		logdata.Message = "Refresh token deletion failure"
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	logdata.Message = "Logout successful."
	logdata.Status = "Success"
	logdata.ResponseCode = http.StatusOK
	logdata.Error = ""
	logger.Info(logdata)
	services.ResponseWithMessage(w, http.StatusOK, nil, "User logged out successfully.", logdata.RequestID)
}

func (h *AuthHandler) AuthMe(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middlewares.UserIDKey).(int64)

	logdata := models.LogData{
		RequestID: r.Context().Value(middlewares.RequestIDKey).(string),
		Endpoint:  r.RequestURI,
		Method:    r.Method,
		Operation: "User Me",
		Status:    "Init",
		UserID:    strconv.FormatInt(userID, 10),
		Message:   "User me initiated.",
	}
	logger.Info(logdata)

	userData, err := h.UserRepo.GetUserByID(userID)
	if err != nil {
		logdata.Message = "User fetch failure"
		logdata.Status = "Error"
		logdata.ResponseCode = http.StatusInternalServerError
		logdata.Error = err.Error()
		logger.Error(logdata)
		services.ResponseWithMessage(w, http.StatusInternalServerError, nil, "Something went wrong.", logdata.RequestID)
		return
	}

	meRes := models.MeResponse{
		ID:        userData.ID,
		UUID:      userData.UUID,
		Name:      userData.Name,
		Email:     userData.Email,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
	}

	logdata.Message = "User me successful."
	logdata.Status = "Success"
	logdata.ResponseCode = http.StatusOK
	logdata.Error = ""
	logger.Info(logdata)
	services.ResponseWithData(w, http.StatusOK, nil, "User found successfully.", meRes, logdata.RequestID)
}

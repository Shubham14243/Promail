package routes

import (
	"net/http"
	"promail/handlers"
	"promail/middlewares"
)

func AuthRoutes(mux *http.ServeMux, h *handlers.AuthHandler) {

	mux.Handle("GET /api/v1/auth/me", middlewares.Auth(http.HandlerFunc(h.AuthMe)))

	mux.HandleFunc("POST /api/v1/auth/signup", h.SignUp)
	mux.HandleFunc("POST /api/v1/auth/login", h.Login)
	mux.HandleFunc("POST /api/v1/auth/refresh", h.RefreshToken)
	mux.Handle("POST /api/v1/auth/logout", middlewares.Auth(http.HandlerFunc(h.Logout)))

}

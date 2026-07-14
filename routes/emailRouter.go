package routes

import (
	"net/http"
	"promail/handlers"
	"promail/middlewares"
)

func EmailRoutes(mux *http.ServeMux, h *handlers.EmailHandler) {

	mux.Handle("POST /api/v1/email/send", middlewares.Auth(http.HandlerFunc(h.SendEmail)))
}

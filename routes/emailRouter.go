package routes

import (
	"net/http"
	"promail/handlers"
	"promail/middlewares"
)

func EmailRoutes(mux *http.ServeMux, h *handlers.EmailHandler) {

	mux.Handle("GET /api/v1/email/track/open/{token}", http.HandlerFunc(h.OpenTrack))
	mux.Handle("GET /api/v1/email/track/click/{token}", http.HandlerFunc(h.ClickTrack))
	mux.Handle("GET /api/v1/email/logs/{uuid}", middlewares.Auth(http.HandlerFunc(h.EmailLogUUID)))
	mux.Handle("GET /api/v1/email/logs", middlewares.Auth(http.HandlerFunc(h.EmailLogAll)))

	mux.Handle("POST /api/v1/email/send", middlewares.Auth(http.HandlerFunc(h.SendEmail)))
	mux.Handle("POST /api/v1/email/send/test", middlewares.Auth(http.HandlerFunc(h.SendEmailTest)))
}

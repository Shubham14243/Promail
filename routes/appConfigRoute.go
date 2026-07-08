package routes

import (
	"net/http"
	"promail/handlers"
	"promail/middlewares"
)

func AppConfigRoutes(mux *http.ServeMux, h *handlers.AppConfigHandler) {

	mux.Handle("GET /api/v1/config/{appID}", middlewares.Auth(http.HandlerFunc(h.GetAppConfigData)))

	mux.Handle("POST /api/v1/config", middlewares.Auth(http.HandlerFunc(h.CreateAppConfig)))
	mux.Handle("PUT /api/v1/config/{appID}", middlewares.Auth(http.HandlerFunc(h.UpdateAppConfig)))
}

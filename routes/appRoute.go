package routes

import (
	"net/http"
	"promail/handlers"
	"promail/middlewares"
)

func AppRoutes(mux *http.ServeMux, h *handlers.AppHandler) {

	mux.Handle("GET /api/v1/apps", middlewares.Auth(http.HandlerFunc(h.GetUserApps)))
	mux.Handle("GET /api/v1/apps/{appID}", middlewares.Auth(http.HandlerFunc(h.GetAppSingle)))

	mux.Handle("GET /api/v1/apps/{appID}/key", middlewares.Auth(http.HandlerFunc(h.GetAppKey)))

	mux.Handle("POST /api/v1/apps", middlewares.Auth(http.HandlerFunc(h.CreateApp)))
	mux.Handle("PUT /api/v1/apps/{appID}", middlewares.Auth(http.HandlerFunc(h.UpdateApp)))
	mux.Handle("DELETE /api/v1/apps/{appID}", middlewares.Auth(http.HandlerFunc(h.DeleteApp)))

	mux.Handle("PUT /api/v1/apps/{appID}/key", middlewares.Auth(http.HandlerFunc(h.UpdateKey)))

}

package routes

import (
	"net/http"
	"promail/handlers"
	"promail/middlewares"
)

func TemplateRoutes(mux *http.ServeMux, h *handlers.TemplateHandler) {

	mux.Handle("GET /api/v1/apps/{appID}/templates", middlewares.Auth(http.HandlerFunc(h.GetAppTemplates)))
	mux.Handle("GET /api/v1/templates/{templateID}", middlewares.Auth(http.HandlerFunc(h.GetTemplateData)))

	mux.Handle("POST /api/v1/templates", middlewares.Auth(http.HandlerFunc(h.CreateTemplate)))
	mux.Handle("PUT /api/v1/templates/{templateID}", middlewares.Auth(http.HandlerFunc(h.UpdateTemplate)))
	mux.Handle("PUT /api/v1/templates/{templateID}/content", middlewares.Auth(http.HandlerFunc(h.UpdateContent)))
	mux.Handle("DELETE /api/v1/templates/{templateID}", middlewares.Auth(http.HandlerFunc(h.DeleteTemplate)))

}

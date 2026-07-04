package routes

import (
	"net/http"
	"promail/handlers"
	"promail/middlewares"
)

func UserRoutes(mux *http.ServeMux, h *handlers.UserHandler) {

	mux.Handle("GET /api/v1/users/all", middlewares.Auth(http.HandlerFunc(h.GetUsers)))
	mux.Handle("GET /api/v1/users", middlewares.Auth(http.HandlerFunc(h.GetUser)))

	mux.Handle("PUT /api/v1/users", middlewares.Auth(http.HandlerFunc(h.UpdateUser)))
	mux.Handle("DELETE /api/v1/users", middlewares.Auth(http.HandlerFunc(h.DeleteUser)))

}

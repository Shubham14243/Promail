package routes

import (
	"net/http"
	"promail/handlers"
	"promail/middlewares"
)

func UserRoutes(mux *http.ServeMux, h *handlers.UserHandler) {

	mux.Handle("PUT /api/v1/users", middlewares.Auth(http.HandlerFunc(h.UpdateUser)))
	mux.Handle("PUT /api/v1/users/password", middlewares.Auth(http.HandlerFunc(h.UpdatePassword)))
	mux.Handle("DELETE /api/v1/users", middlewares.Auth(http.HandlerFunc(h.DeleteUser)))

}

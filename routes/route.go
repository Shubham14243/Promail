package routes

import (
	"net/http"
	"promail/handlers"
)

func RegisterRoutes(uh *handlers.UserHandler, ah *handlers.AuthHandler, aph *handlers.AppHandler, th *handlers.TemplateHandler) *http.ServeMux {

	mux := http.NewServeMux()

	AuthRoutes(mux, ah)
	UserRoutes(mux, uh)
	AppRoutes(mux, aph)
	TemplateRoutes(mux, th)

	return mux
}

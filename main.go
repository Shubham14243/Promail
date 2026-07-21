package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"promail/configs"
	"promail/handlers"
	"promail/logger"
	"promail/middlewares"
	"promail/repositories"
	"promail/routes"
	"promail/services"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]string{
		"ping": "pong",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {

	response := map[string]string{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func main() {

	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/health", healthHandler)

	configs.LoadEnv()
	log.Println("Connecting DB...")
	configs.ConnectDB()
	log.Println("DB Connected")

	log.Println("Running migrations...")
	if err := configs.Migrate(); err != nil {
		log.Fatal(err)
	}
	log.Println("Migration completed")
	logger.Init()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	emailWorker := services.NewWorker(configs.DB)
	go emailWorker.Run(ctx)

	userRepo := &repositories.UserRepository{
		DB: configs.DB,
	}
	refreshTokenRepo := &repositories.RefreshTokenRepository{
		DB: configs.DB,
	}
	appRepo := &repositories.AppRepository{
		DB: configs.DB,
	}
	tempRepo := &repositories.TemplateRepository{
		DB: configs.DB,
	}
	appConfigRepo := &repositories.AppConfigRepository{
		DB: configs.DB,
	}
	emailRepo := &repositories.EmailRepository{
		DB: configs.DB,
	}

	UserHandler := &handlers.UserHandler{
		Repo: userRepo,
	}
	AuthHandler := &handlers.AuthHandler{
		UserRepo:         userRepo,
		RefreshTokenRepo: refreshTokenRepo,
	}
	AppHandler := &handlers.AppHandler{
		AppRepo: appRepo,
	}
	TemplateHandler := &handlers.TemplateHandler{
		TempRepo: tempRepo,
	}
	AppConfigHandler := &handlers.AppConfigHandler{
		AppConfigRepo: appConfigRepo,
	}
	EmailHandler := &handlers.EmailHandler{
		EmailRepo:     emailRepo,
		AppConfigRepo: appConfigRepo,
		AppRepo:       appRepo,
		TempRepo:      tempRepo,
	}

	mux := routes.RegisterRoutes(UserHandler, AuthHandler, AppHandler, TemplateHandler, AppConfigHandler, EmailHandler)

	mux.HandleFunc("GET /ping", pingHandler)
	mux.HandleFunc("GET /health", healthHandler)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port

	log.Println("Starting server on", addr)

	routeHandler := middlewares.RequestID(mux)

	err := http.ListenAndServe(addr, routeHandler)
	if err != nil {
		log.Fatal(err)
	}

}

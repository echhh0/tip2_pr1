package main

import (
	"log"
	"net/http"
	"os"

	httpapi "tip2_pr1/services/auth/internal/http"
	"tip2_pr1/services/auth/internal/service"
	"tip2_pr1/shared/middleware"
)

func main() {
	port := getEnv("AUTH_PORT", "8081")

	authService := service.New()
	handler := httpapi.New(authService)

	mux := http.NewServeMux()
	handler.Register(mux)

	app := middleware.RequestID(middleware.Logging(mux))

	addr := ":" + port
	log.Printf("auth service starting on %s", addr)

	if err := http.ListenAndServe(addr, app); err != nil {
		log.Fatalf("auth service failed: %v", err)
	}
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

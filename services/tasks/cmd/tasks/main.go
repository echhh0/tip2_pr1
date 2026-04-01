package main

import (
	"log"
	"net/http"
	"os"

	"tip2_pr1/services/tasks/internal/client/authclient"
	httpapi "tip2_pr1/services/tasks/internal/http"
	"tip2_pr1/services/tasks/internal/service"
	"tip2_pr1/shared/middleware"
)

func main() {
	port := getEnv("TASKS_PORT", "8082")
	authBaseURL := getEnv("AUTH_BASE_URL", "http://localhost:8081")

	taskService := service.New()
	authClient := authclient.New(authBaseURL)
	handler := httpapi.New(taskService, authClient)

	mux := http.NewServeMux()
	handler.Register(mux)

	app := middleware.RequestID(middleware.Logging(mux))

	addr := ":" + port
	log.Printf("tasks service starting on %s, auth_base_url=%s", addr, authBaseURL)

	if err := http.ListenAndServe(addr, app); err != nil {
		log.Fatalf("tasks service failed: %v", err)
	}
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

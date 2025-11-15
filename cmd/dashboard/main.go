package main

import (
	"crypto/subtle"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/phoenix-marie/core/internal/api"
)

// Basic security middleware
func basicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for static files and websocket connections
		if r.URL.Path == "/ws" || r.URL.Path == "/" || r.URL.Path == "/css/styles.css" || r.URL.Path == "/js/app.js" {
			next.ServeHTTP(w, r)
			return
		}

		// Get API key from header
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Compare API key (in production, use environment variables or secure configuration)
		expectedKey := "phoenix-dashboard-key"
		if subtle.ConstantTimeCompare([]byte(apiKey), []byte(expectedKey)) != 1 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Initialize server and metrics service
	server := api.NewServer()
	metricsService := api.NewMetricsService(server)

	// Start services
	server.Start()
	metricsService.Start()

	// Set up HTTP server with security middleware
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: basicAuth(server.SetupRoutes()),
	}

	// Handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Dashboard server starting on http://localhost:8080")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-stop
	log.Println("Shutting down server...")
}

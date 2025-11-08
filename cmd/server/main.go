// Main application entry point with proper initialization and graceful shutdown
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/domain-reference-Master-Geopolitical/internal/xcut/bootstrap"
	"github.com/domain-reference-Master-Geopolitical/internal/xcut/logging"
)

func main() {
	ctx := context.Background()

	// Initialize application container
	container, err := bootstrap.InitializeContainer(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", container.Config.Server.Host, container.Config.Server.Port),
		ReadTimeout:  container.Config.Server.ReadTimeout,
		WriteTimeout: container.Config.Server.WriteTimeout,
		Handler:      setupRoutes(container),
	}

	// Start server in goroutine
	go func() {
		container.Logger.Info(ctx, "Starting HTTP server",
			logging.Field{Key: "address", Value: server.Addr},
		)
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			container.Logger.Error(ctx, "Server failed to start", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	container.Logger.Info(ctx, "Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		container.Logger.Error(ctx, "Server forced to shutdown", err)
	}

	// Cleanup container
	if err := container.Shutdown(ctx); err != nil {
		container.Logger.Error(ctx, "Error during container shutdown", err)
	}

	container.Logger.Info(ctx, "Server exited")
}

func setupRoutes(container *bootstrap.Container) http.Handler {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"geopolitical-service"}`))
	})

	// Ready check endpoint
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ready","service":"geopolitical-service"}`))
	})

	// API routes will be added here
	mux.HandleFunc("/api/v1/countries", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Countries API endpoint - implementation in progress"}`))
	})

	return mux
}
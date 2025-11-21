package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"unsri-backend/internal/api-gateway/config"
	"unsri-backend/internal/api-gateway/handler"
	"unsri-backend/internal/shared/logger"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	log := logger.New(cfg.LogLevel)
	log.Info("Starting API Gateway...")

	// Setup router
	router := gin.Default()
	router.Use(gin.Recovery())

	// Initialize proxy handler
	proxyHandler := handler.NewProxyHandler(cfg, log)

	// Setup routes
	handler.SetupRoutes(router, proxyHandler)

	// Setup Swagger (only in development)
	if cfg.LogLevel == "debug" || os.Getenv("ENABLE_SWAGGER") == "true" {
		setupSwagger(router)
		log.Info("Swagger UI available at /swagger/index.html")
	}

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", err)
		}
	}()

	log.Infof("API Gateway started on port %s", cfg.Port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", err)
	}

	log.Info("Server exited")
}


package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"unsri-backend/internal/api-gateway/config"
	"unsri-backend/internal/api-gateway/handler"
	"unsri-backend/internal/api-gateway/service"
	"unsri-backend/internal/shared/logger"
	"unsri-backend/internal/shared/messaging"

	"github.com/gin-gonic/gin"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	log := logger.New(cfg.LogLevel)
	log.Info("Starting API Gateway...")

	// Initialize RabbitMQ connection with retry logic
	var rabbitMQClient *messaging.RabbitMQClient
	maxRetries := 5
	retryDelay := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		var err error
		rabbitMQClient, err = messaging.NewRabbitMQ(messaging.Config{
			Host:     cfg.RabbitMQHost,
			Port:     cfg.RabbitMQPort,
			User:     cfg.RabbitMQUser,
			Password: cfg.RabbitMQPassword,
			VHost:    cfg.RabbitMQVHost,
		})
		if err == nil {
			log.Info("Connected to RabbitMQ")
			break
		}

		if i < maxRetries-1 {
			log.Warnf("Failed to connect to RabbitMQ (attempt %d/%d): %v. Retrying in %v...", i+1, maxRetries, err, retryDelay)
			time.Sleep(retryDelay)
		} else {
			log.Fatalf("Failed to connect to RabbitMQ after %d attempts: %v", maxRetries, err)
		}
	}
	defer rabbitMQClient.Close()

	// Initialize message broker service
	messageBrokerService := service.NewMessageBrokerService(rabbitMQClient, log)
	if err := messageBrokerService.Initialize(); err != nil {
		log.Fatalf("Failed to initialize message broker: %v", err)
	}
	log.Info("Message broker initialized")

	// Start consuming messages (optional - for handling responses from services)
	// This can be used for async operations or event-driven workflows
	go func() {
		if err := messageBrokerService.StartConsumer(
			"notification_queue",
			"api_gateway_consumer",
			func(msg amqp.Delivery) error {
				log.Infof("Received notification message: %s", string(msg.Body))
				// Handle notification message
				// You can process notifications, forward to other services, etc.
				return nil
			},
		); err != nil {
			log.Errorf("Failed to start consumer: %v", err)
		}
	}()

	// Setup router
	router := gin.Default()
	router.Use(gin.Recovery())

	// Initialize proxy handler with message broker
	proxyHandler := handler.NewProxyHandler(cfg, log, messageBrokerService)

	// Setup routes
	handler.SetupRoutes(router, proxyHandler)

	// Setup Swagger (only in development)
	// Uncomment if swagger is needed
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

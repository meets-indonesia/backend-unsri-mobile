package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"unsri-backend/internal/auth/config"
	"unsri-backend/internal/auth/handler"
	"unsri-backend/internal/auth/repository"
	"unsri-backend/internal/auth/service"
	"unsri-backend/internal/shared/database"
	"unsri-backend/internal/shared/logger"
	"unsri-backend/internal/shared/models"
	"unsri-backend/pkg/jwt"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	log := logger.New(cfg.LogLevel)
	log.Info("Starting auth service...")

	// Initialize database
	db, err := database.NewPostgres(database.Config{
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		User:            cfg.Database.User,
		Password:        cfg.Database.Password,
		DBName:          cfg.Database.DBName,
		SSLMode:         cfg.Database.SSLMode,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	})
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	// Auto migrate models
	if err := db.AutoMigrate(
		&models.User{},
		&models.Mahasiswa{},
		&models.Dosen{},
		&models.Staff{},
	); err != nil {
		log.Fatal("Failed to migrate database", err)
	}

	// Initialize JWT
	jwtToken := jwt.NewJWT(
		cfg.JWT.SecretKey,
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
	)

	// Initialize repository
	authRepo := repository.NewAuthRepository(db)

	// Initialize service
	authService := service.NewAuthService(authRepo, jwtToken)

	// Initialize handler
	authHandler := handler.NewAuthHandler(authService, log)

	// Setup router
	router := gin.Default()
	router.Use(gin.Recovery())
	handler.SetupRoutes(router, authHandler)

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

	log.Infof("Auth service started on port %s", cfg.Port)

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


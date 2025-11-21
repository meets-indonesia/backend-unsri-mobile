package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"unsri-backend/internal/attendance/config"
	"unsri-backend/internal/attendance/handler"
	"unsri-backend/internal/attendance/repository"
	"unsri-backend/internal/attendance/service"
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
	log.Info("Starting attendance service...")

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
		&models.Attendance{},
		&models.AttendanceSession{},
		&models.Schedule{},
	); err != nil {
		log.Fatal("Failed to migrate database", err)
	}

	// Initialize JWT
	jwtToken := jwt.NewJWT(
		cfg.JWT.SecretKey,
		15*time.Minute,  // Access token TTL
		7*24*time.Hour,  // Refresh token TTL
	)

	// Initialize repository
	attendanceRepo := repository.NewAttendanceRepository(db)

	// Initialize service
	attendanceService := service.NewAttendanceService(attendanceRepo, jwtToken)

	// Initialize handler
	attendanceHandler := handler.NewAttendanceHandler(attendanceService, log)

	// Setup router
	router := gin.Default()
	router.Use(gin.Recovery())
	handler.SetupRoutes(router, attendanceHandler, jwtToken)

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

	log.Infof("Attendance service started on port %s", cfg.Port)

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


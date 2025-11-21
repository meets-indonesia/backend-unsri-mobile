package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"unsri-backend/internal/shared/database"
	"unsri-backend/internal/shared/logger"
	"unsri-backend/internal/shared/models"
	"unsri-backend/internal/schedule/config"
	"unsri-backend/internal/schedule/handler"
	"unsri-backend/internal/schedule/repository"
	"unsri-backend/internal/schedule/service"
	"unsri-backend/pkg/jwt"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel)
	log.Info("Starting schedule service...")

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

	if err := db.AutoMigrate(
		&models.Schedule{},
	); err != nil {
		log.Fatal("Failed to migrate database", err)
	}

	jwtToken := jwt.NewJWT(
		cfg.JWT.SecretKey,
		15*time.Minute,
		7*24*time.Hour,
	)

	scheduleRepo := repository.NewScheduleRepository(db)
	scheduleService := service.NewScheduleService(scheduleRepo)
	scheduleHandler := handler.NewScheduleHandler(scheduleService, log)

	router := gin.Default()
	router.Use(gin.Recovery())
	handler.SetupRoutes(router, scheduleHandler, jwtToken)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", err)
		}
	}()

	log.Infof("Schedule service started on port %s", cfg.Port)

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


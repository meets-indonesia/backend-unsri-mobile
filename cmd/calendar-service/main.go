package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"unsri-backend/internal/calendar/config"
	"unsri-backend/internal/calendar/handler"
	"unsri-backend/internal/calendar/repository"
	"unsri-backend/internal/calendar/service"
	"unsri-backend/internal/shared/database"
	"unsri-backend/internal/shared/logger"
	"unsri-backend/internal/shared/models"
	"unsri-backend/pkg/jwt"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel)
	log.Info("Starting academic calendar service...")

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
		&models.AcademicEvent{},
	); err != nil {
		log.Fatal("Failed to migrate database", err)
	}

	jwtToken := jwt.NewJWT(
		cfg.JWT.SecretKey,
		15*time.Minute,
		7*24*time.Hour,
	)

	calendarRepo := repository.NewCalendarRepository(db)
	calendarService := service.NewCalendarService(calendarRepo)
	calendarHandler := handler.NewCalendarHandler(calendarService, log)

	router := gin.Default()
	router.Use(gin.Recovery())
	handler.SetupRoutes(router, calendarHandler, jwtToken)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", err)
		}
	}()

	log.Infof("Academic calendar service started on port %s", cfg.Port)

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


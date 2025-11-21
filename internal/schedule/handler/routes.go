package handler

import (
	"github.com/gin-gonic/gin"
	"unsri-backend/internal/schedule/middleware"
	"unsri-backend/pkg/jwt"
)

// SetupRoutes sets up all routes for schedule service
func SetupRoutes(router *gin.Engine, handler *ScheduleHandler, jwtToken *jwt.JWT) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "schedule-service"})
	})

	v1 := router.Group("/api/v1/schedules")
	v1.Use(middleware.AuthMiddleware(jwtToken))
	{
		v1.GET("", handler.GetSchedules)
		v1.GET("/today", handler.GetTodaySchedules)
		v1.GET("/upcoming", handler.GetUpcomingSchedules)
		v1.GET("/calendar/:year/:month", handler.GetCalendarView)
		v1.GET("/:id", handler.GetSchedule)
		v1.POST("", middleware.RoleMiddleware("dosen", "staff"), handler.CreateSchedule)
		v1.PUT("/:id", middleware.RoleMiddleware("dosen", "staff"), handler.UpdateSchedule)
		v1.DELETE("/:id", middleware.RoleMiddleware("dosen", "staff"), handler.DeleteSchedule)
	}
}


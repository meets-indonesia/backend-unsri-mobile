package handler

import (
	"github.com/gin-gonic/gin"
	"unsri-backend/internal/report/middleware"
	"unsri-backend/pkg/jwt"
)

// SetupRoutes sets up all routes for report service
func SetupRoutes(router *gin.Engine, handler *ReportHandler, jwtToken *jwt.JWT) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "report-service"})
	})

	v1 := router.Group("/api/v1/reports")
	v1.Use(middleware.AuthMiddleware(jwtToken))
	{
		v1.GET("/attendance", handler.GetAttendanceReport)
		v1.GET("/academic", handler.GetAcademicReport)
		v1.GET("/course", handler.GetCourseReport)
		v1.GET("/daily", handler.GetDailyReport)
	}
}


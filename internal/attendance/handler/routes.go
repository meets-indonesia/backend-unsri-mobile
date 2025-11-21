package handler

import (
	"github.com/gin-gonic/gin"
	"unsri-backend/internal/attendance/middleware"
	"unsri-backend/pkg/jwt"
)

// SetupRoutes sets up all routes for attendance service
func SetupRoutes(router *gin.Engine, handler *AttendanceHandler, jwtToken *jwt.JWT) {
	// Public routes
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "attendance-service"})
	})

	// Protected routes
	v1 := router.Group("/api/v1/attendance")
	v1.Use(middleware.AuthMiddleware(jwtToken))
	{
		// QR code operations
		v1.POST("/qr/generate", middleware.RoleMiddleware("dosen", "staff"), handler.GenerateQR)
		v1.POST("/qr/scan", handler.ScanQR)

		// Attendance operations
		v1.GET("", handler.GetAttendances)
		v1.GET("/statistics", handler.GetStatistics)
		v1.GET("/overview", handler.GetOverview)
		v1.GET("/history", handler.GetAttendances) // Alias for GetAttendances
		v1.GET("/by-course/:courseId", handler.GetByCourse)
		v1.GET("/by-student/:studentId", handler.GetByStudent)
		v1.POST("/manual", middleware.RoleMiddleware("dosen", "staff"), handler.CreateManualAttendance)
		v1.PUT("/:id", middleware.RoleMiddleware("dosen", "staff"), handler.UpdateAttendance)

		// Campus attendance (tap in/out)
		v1.POST("/tap-in", handler.TapIn)
		v1.POST("/tap-out", handler.TapOut)
	}

	// Schedule routes
	schedules := router.Group("/api/v1/schedules")
	schedules.Use(middleware.AuthMiddleware(jwtToken))
	{
		schedules.GET("", handler.GetSchedules)
		schedules.GET("/today", handler.GetTodaySchedules)
		schedules.GET("/:id", handler.GetSchedule)
		schedules.POST("", middleware.RoleMiddleware("dosen", "staff"), handler.CreateSchedule)
		schedules.PUT("/:id", middleware.RoleMiddleware("dosen", "staff"), handler.UpdateSchedule)
		schedules.DELETE("/:id", middleware.RoleMiddleware("dosen", "staff"), handler.DeleteSchedule)
	}
}


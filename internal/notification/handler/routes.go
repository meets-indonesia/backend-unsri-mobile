package handler

import (
	"github.com/gin-gonic/gin"
	"unsri-backend/internal/notification/middleware"
	"unsri-backend/pkg/jwt"
)

// SetupRoutes sets up all routes for notification service
func SetupRoutes(router *gin.Engine, handler *NotificationHandler, jwtToken *jwt.JWT) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "notification-service"})
	})

	v1 := router.Group("/api/v1/notifications")
	v1.Use(middleware.AuthMiddleware(jwtToken))
	{
		v1.POST("/send", middleware.RoleMiddleware("dosen", "staff"), handler.SendNotification)
		v1.GET("", handler.GetNotifications)
		v1.PUT("/:id/read", handler.MarkAsRead)
		v1.PUT("/read-all", handler.MarkAllAsRead)
		v1.POST("/register-device", handler.RegisterDeviceToken)
		v1.DELETE("/device/:token", handler.UnregisterDeviceToken)
	}
}


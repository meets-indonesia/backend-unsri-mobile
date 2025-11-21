package handler

import (
	"github.com/gin-gonic/gin"
	"unsri-backend/internal/quick-actions/middleware"
	"unsri-backend/pkg/jwt"
)

// SetupRoutes sets up all routes for quick actions service
func SetupRoutes(router *gin.Engine, handler *QuickActionsHandler, jwtToken *jwt.JWT) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "quick-actions-service"})
	})

	v1 := router.Group("/api/v1/quick-actions")
	v1.Use(middleware.AuthMiddleware(jwtToken))
	{
		v1.GET("", handler.GetQuickActions)
		v1.GET("/transcript/:studentId", handler.GetTranscript)
		v1.GET("/krs/:studentId", handler.GetKRS)
		v1.GET("/bimbingan", handler.GetBimbingans)
	}
}


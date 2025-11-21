package handler

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up all routes for auth service
func SetupRoutes(router *gin.Engine, handler *AuthHandler) {
	v1 := router.Group("/api/v1/auth")
	{
		v1.POST("/login", handler.Login)
		v1.POST("/register", handler.Register)
		v1.POST("/refresh", handler.RefreshToken)
		v1.GET("/verify", handler.VerifyToken)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "auth-service"})
	})
}


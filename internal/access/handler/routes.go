package handler

import (
	"github.com/gin-gonic/gin"
	"unsri-backend/internal/access/middleware"
	"unsri-backend/pkg/jwt"
)

// SetupRoutes sets up all routes for access service
func SetupRoutes(router *gin.Engine, handler *AccessHandler, jwtToken *jwt.JWT) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "access-service"})
	})

	v1 := router.Group("/api/v1/access")
	v1.Use(middleware.AuthMiddleware(jwtToken))
	{
		v1.POST("/qr/validate", handler.ValidateQR)
		v1.GET("/history", handler.GetAccessHistory)
		v1.POST("/log", handler.LogAccess)
		v1.GET("/permissions/:userId", handler.GetAccessPermissions)
		v1.POST("/permissions", middleware.RoleMiddleware("staff"), handler.CreateAccessPermission)
	}
}


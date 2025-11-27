package handler

import (
	"github.com/gin-gonic/gin"
	"unsri-backend/internal/qr/middleware"
	"unsri-backend/pkg/jwt"
)

// SetupRoutes sets up all routes for QR service
func SetupRoutes(router *gin.Engine, handler *QRHandler, jwtToken *jwt.JWT) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "qr-service"})
	})

	// Public endpoint for gate UNSRI validation (no auth required)
	public := router.Group("/api/v1/qr")
	{
		// Gate UNSRI validation endpoint (public)
		public.POST("/gate/validate", handler.ValidateGateQR)
	}

	v1 := router.Group("/api/v1/qr")
	v1.Use(middleware.AuthMiddleware(jwtToken))
	{
		// General QR operations
		v1.POST("/generate", handler.GenerateQR)
		v1.POST("/validate", handler.ValidateQR)
		v1.GET("/:id", handler.GetQR)

		// Class attendance QR (regenerates after each scan)
		v1.POST("/class/generate", middleware.RoleMiddleware("dosen", "staff"), handler.GenerateClassQR)
		v1.POST("/class/:scheduleId/regenerate", middleware.RoleMiddleware("dosen", "staff"), handler.RegenerateClassQR)

		// Gate access QR (unique per user, doesn't change)
		v1.GET("/access/generate", handler.GenerateAccessQR)
		v1.GET("/access/validate/:token", handler.ValidateAccessQR)
	}
}


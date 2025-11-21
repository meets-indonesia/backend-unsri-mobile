package handler

import (
	"github.com/gin-gonic/gin"
	"unsri-backend/internal/location/middleware"
	"unsri-backend/pkg/jwt"
)

// SetupRoutes sets up all routes for location service
func SetupRoutes(router *gin.Engine, handler *LocationHandler, jwtToken *jwt.JWT) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "location-service"})
	})

	v1 := router.Group("/api/v1/location")
	v1.Use(middleware.AuthMiddleware(jwtToken))
	{
		v1.POST("/tap-in", handler.TapIn)
		v1.POST("/tap-out", handler.TapOut)
		v1.GET("/check-in-status", handler.GetCheckInStatus)
		v1.GET("/history", handler.GetLocationHistory)
		v1.GET("/geofences", handler.GetGeofences)
		v1.POST("/validate", handler.ValidateLocation)
		v1.POST("/geofences", middleware.RoleMiddleware("staff"), handler.CreateGeofence)
	}
}


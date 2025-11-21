package handler

import (
	"github.com/gin-gonic/gin"
	"unsri-backend/internal/broadcast/middleware"
	"unsri-backend/pkg/jwt"
)

// SetupRoutes sets up all routes for broadcast service
func SetupRoutes(router *gin.Engine, handler *BroadcastHandler, jwtToken *jwt.JWT) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "broadcast-service"})
	})

	v1 := router.Group("/api/v1/broadcasts")
	v1.Use(middleware.AuthMiddleware(jwtToken))
	{
		v1.GET("", handler.GetBroadcasts)
		v1.GET("/general", handler.GetGeneralBroadcasts)
		v1.GET("/class", handler.GetClassBroadcasts)
		v1.GET("/:id", handler.GetBroadcast)
		v1.POST("", middleware.RoleMiddleware("dosen", "staff"), handler.CreateBroadcast)
		v1.PUT("/:id", middleware.RoleMiddleware("dosen", "staff"), handler.UpdateBroadcast)
		v1.DELETE("/:id", middleware.RoleMiddleware("dosen", "staff"), handler.DeleteBroadcast)
		v1.POST("/:id/schedule", middleware.RoleMiddleware("dosen", "staff"), handler.ScheduleBroadcast)
		v1.POST("/search", handler.SearchBroadcasts)
	}
}


package handler

import (
	"github.com/gin-gonic/gin"
	"unsri-backend/internal/search/middleware"
	"unsri-backend/pkg/jwt"
)

// SetupRoutes sets up all routes for search service
func SetupRoutes(router *gin.Engine, handler *SearchHandler, jwtToken *jwt.JWT) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "search-service"})
	})

	v1 := router.Group("/api/v1/search")
	v1.Use(middleware.AuthMiddleware(jwtToken))
	{
		v1.GET("", handler.Search)
		v1.GET("/global", handler.GlobalSearch)
	}
}


package handler

import (
	"github.com/gin-gonic/gin"
	"unsri-backend/internal/calendar/middleware"
	"unsri-backend/pkg/jwt"
)

// SetupRoutes sets up all routes for calendar service
func SetupRoutes(router *gin.Engine, handler *CalendarHandler, jwtToken *jwt.JWT) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "academic-calendar-service"})
	})

	v1 := router.Group("/api/v1/calendar/events")
	v1.Use(middleware.AuthMiddleware(jwtToken))
	{
		v1.GET("", handler.GetEvents)
		v1.GET("/upcoming", handler.GetUpcomingEvents)
		v1.GET("/month/:year/:month", handler.GetEventsByMonth)
		v1.GET("/:id", handler.GetEvent)
		v1.POST("", middleware.RoleMiddleware("staff"), handler.CreateEvent)
		v1.PUT("/:id", middleware.RoleMiddleware("staff"), handler.UpdateEvent)
		v1.DELETE("/:id", middleware.RoleMiddleware("staff"), handler.DeleteEvent)
	}
}


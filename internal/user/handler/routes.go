package handler

import (
	"github.com/gin-gonic/gin"
	"unsri-backend/internal/user/middleware"
	"unsri-backend/pkg/jwt"
)

// SetupRoutes sets up all routes for user service
func SetupRoutes(router *gin.Engine, handler *UserHandler, jwtToken *jwt.JWT) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "user-service"})
	})

	// Protected routes
	v1 := router.Group("/api/v1/users")
	v1.Use(middleware.AuthMiddleware(jwtToken))
	{
		v1.GET("/profile", handler.GetProfile)
		v1.PUT("/profile", handler.UpdateProfile)
		v1.POST("/avatar", handler.UploadAvatar)
		v1.GET("/search", handler.SearchUsers)
		v1.GET("/:id", handler.GetUserByID)
		v1.GET("/mahasiswa/:nim", handler.GetMahasiswaByNIM)
		v1.GET("/dosen/:nip", handler.GetDosenByNIP)
		v1.GET("/staff/:nip", handler.GetStaffByNIP)
	}
}


package handler

import (
	"unsri-backend/internal/user/middleware"
	"unsri-backend/pkg/jwt"

	"github.com/gin-gonic/gin"
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
		// User profile (self)
		v1.GET("/profile", handler.GetProfile)
		v1.PUT("/profile", handler.UpdateProfile)
		v1.POST("/avatar", handler.UploadAvatar)

		// Search users
		v1.GET("/search", handler.SearchUsers)

		// Get user by identifier (specific routes first to avoid conflict)
		v1.GET("/mahasiswa/:nim", handler.GetMahasiswaByNIM)
		v1.GET("/dosen/:nip", handler.GetDosenByNIP)
		v1.GET("/staff/:nip", handler.GetStaffByNIP)

		// Admin only routes (staff role required)
		admin := v1.Group("")
		admin.Use(middleware.RoleMiddleware("staff"))
		{
			// List all users (must be before /:id route)
			admin.GET("", handler.ListUsers)

			// Create user (admin)
			admin.POST("", handler.CreateUser)

			// Update user (admin)
			admin.PUT("/:id", handler.AdminUpdateUser)

			// Delete user (admin)
			admin.DELETE("/:id", handler.DeleteUser)

			// Activate/Deactivate user (admin) - must be before /:id route
			admin.PUT("/:id/activate", handler.ActivateUser)
			admin.PUT("/:id/deactivate", handler.DeactivateUser)
		}

		// Get user by ID (must be last to avoid conflict with other routes)
		v1.GET("/:id", handler.GetUserByID)
	}
}

package handler

import (
	"github.com/gin-gonic/gin"
	"unsri-backend/internal/file-storage/middleware"
	"unsri-backend/pkg/jwt"
)

// SetupRoutes sets up all routes for file storage service
func SetupRoutes(router *gin.Engine, handler *FileStorageHandler, jwtToken *jwt.JWT) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "file-storage-service"})
	})

	v1 := router.Group("/api/v1/files")
	v1.Use(middleware.AuthMiddleware(jwtToken))
	{
		v1.POST("/upload", handler.UploadFile)
		v1.GET("", handler.GetFiles)
		v1.GET("/:id", handler.GetFile)
		v1.GET("/:id/download", handler.DownloadFile)
		v1.DELETE("/:id", handler.DeleteFile)
		v1.POST("/avatar", handler.UploadAvatar)
		v1.POST("/document", handler.UploadDocument)
	}
}


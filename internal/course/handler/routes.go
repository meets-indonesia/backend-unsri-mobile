package handler

import (
	"github.com/gin-gonic/gin"
	"unsri-backend/internal/course/middleware"
	"unsri-backend/pkg/jwt"
)

// SetupRoutes sets up all routes for course service
func SetupRoutes(router *gin.Engine, handler *CourseHandler, jwtToken *jwt.JWT) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "course-service"})
	})

	v1 := router.Group("/api/v1/courses")
	v1.Use(middleware.AuthMiddleware(jwtToken))
	{
		v1.GET("", handler.GetCourses)
		v1.GET("/:id", handler.GetCourse)
		v1.POST("", middleware.RoleMiddleware("dosen", "staff"), handler.CreateCourse)
		v1.PUT("/:id", middleware.RoleMiddleware("dosen", "staff"), handler.UpdateCourse)
		v1.DELETE("/:id", middleware.RoleMiddleware("dosen", "staff"), handler.DeleteCourse)
		v1.GET("/by-student/:studentId", handler.GetClassesByStudent)
		v1.GET("/by-lecturer/:lecturerId", handler.GetClassesByLecturer)
	}

	classes := router.Group("/api/v1/classes")
	classes.Use(middleware.AuthMiddleware(jwtToken))
	{
		classes.GET("", handler.GetClasses)
		classes.GET("/:id", handler.GetClass)
		classes.POST("", middleware.RoleMiddleware("dosen", "staff"), handler.CreateClass)
	}
}


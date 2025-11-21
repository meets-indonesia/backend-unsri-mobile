package handler

import (
	"github.com/gin-gonic/gin"
	"unsri-backend/internal/course/service"
	"unsri-backend/internal/shared/logger"
	"unsri-backend/internal/shared/utils"
)

// CourseHandler handles HTTP requests for course management
type CourseHandler struct {
	service *service.CourseService
	logger  logger.Logger
}

// NewCourseHandler creates a new course handler
func NewCourseHandler(service *service.CourseService, logger logger.Logger) *CourseHandler {
	return &CourseHandler{
		service: service,
		logger:  logger,
	}
}

// CreateCourse handles create course request
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	var req service.CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.CreateCourse(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 201, result)
}

// GetCourse handles get course by ID request
func (h *CourseHandler) GetCourse(c *gin.Context) {
	courseID := c.Param("id")

	result, err := h.service.GetCourseByID(c.Request.Context(), courseID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GetCourses handles get courses request
func (h *CourseHandler) GetCourses(c *gin.Context) {
	var req service.GetCoursesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	courses, total, err := h.service.GetCourses(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}

	utils.PaginatedResponse(c, courses, page, perPage, total)
}

// UpdateCourse handles update course request
func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	courseID := c.Param("id")

	var req service.UpdateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.UpdateCourse(c.Request.Context(), courseID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// DeleteCourse handles delete course request
func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	courseID := c.Param("id")

	err := h.service.DeleteCourse(c.Request.Context(), courseID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, gin.H{"message": "Course deleted successfully"})
}

// CreateClass handles create class request
func (h *CourseHandler) CreateClass(c *gin.Context) {
	var req service.CreateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.CreateClass(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 201, result)
}

// GetClass handles get class by ID request
func (h *CourseHandler) GetClass(c *gin.Context) {
	classID := c.Param("id")

	result, err := h.service.GetClassByID(c.Request.Context(), classID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GetClasses handles get classes request
func (h *CourseHandler) GetClasses(c *gin.Context) {
	var req service.GetClassesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	classes, total, err := h.service.GetClasses(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}

	utils.PaginatedResponse(c, classes, page, perPage, total)
}

// GetClassesByStudent handles get classes by student request
func (h *CourseHandler) GetClassesByStudent(c *gin.Context) {
	studentID := c.Param("studentId")

	result, err := h.service.GetClassesByStudent(c.Request.Context(), studentID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GetClassesByLecturer handles get classes by lecturer request
func (h *CourseHandler) GetClassesByLecturer(c *gin.Context) {
	lecturerID := c.Param("lecturerId")

	result, err := h.service.GetClassesByLecturer(c.Request.Context(), lecturerID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}


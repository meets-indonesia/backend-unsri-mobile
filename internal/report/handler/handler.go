package handler

import (
	"unsri-backend/internal/report/service"
	"unsri-backend/internal/shared/logger"
	"unsri-backend/internal/shared/utils"

	"github.com/gin-gonic/gin"
)

// ReportHandler handles HTTP requests for reports
type ReportHandler struct {
	service *service.ReportService
	logger  logger.Logger
}

// NewReportHandler creates a new report handler
func NewReportHandler(service *service.ReportService, logger logger.Logger) *ReportHandler {
	return &ReportHandler{
		service: service,
		logger:  logger,
	}
}

// GetAttendanceReport handles get attendance report request
func (h *ReportHandler) GetAttendanceReport(c *gin.Context) {
	var req service.AttendanceReportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	// If user is student, only allow their own report
	userRole := c.GetString("user_role")
	userID := c.GetString("user_id")
	if userRole == "mahasiswa" && (req.StudentID == nil || *req.StudentID != userID) {
		req.StudentID = &userID
	}

	result, err := h.service.GetAttendanceReport(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GetAcademicReport handles get academic report request
func (h *ReportHandler) GetAcademicReport(c *gin.Context) {
	var req service.AcademicReportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	// If user is student, only allow their own report
	userRole := c.GetString("user_role")
	userID := c.GetString("user_id")
	if userRole == "mahasiswa" && req.StudentID != userID {
		req.StudentID = userID
	}

	result, err := h.service.GetAcademicReport(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GetCourseReport handles get course report request
func (h *ReportHandler) GetCourseReport(c *gin.Context) {
	var req service.CourseReportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.GetCourseReport(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GetDailyReport handles get daily report request
func (h *ReportHandler) GetDailyReport(c *gin.Context) {
	var req service.DailyReportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.GetDailyReport(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

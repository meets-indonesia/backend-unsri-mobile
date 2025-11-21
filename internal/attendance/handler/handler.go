package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"unsri-backend/internal/attendance/service"
	"unsri-backend/internal/shared/logger"
	"unsri-backend/internal/shared/utils"
)

// AttendanceHandler handles HTTP requests for attendance
type AttendanceHandler struct {
	service *service.AttendanceService
	logger  logger.Logger
}

// NewAttendanceHandler creates a new attendance handler
func NewAttendanceHandler(service *service.AttendanceService, logger logger.Logger) *AttendanceHandler {
	return &AttendanceHandler{
		service: service,
		logger:  logger,
	}
}

// GenerateQR handles QR code generation request
func (h *AttendanceHandler) GenerateQR(c *gin.Context) {
	userID := c.GetString("user_id")

	var req service.GenerateQRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.GenerateQRCode(c.Request.Context(), userID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// ScanQR handles QR code scan request
func (h *AttendanceHandler) ScanQR(c *gin.Context) {
	userID := c.GetString("user_id")

	var req service.ScanQRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.ScanQRCode(c.Request.Context(), userID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GetAttendances handles get attendances request
func (h *AttendanceHandler) GetAttendances(c *gin.Context) {
	userID := c.GetString("user_id")
	userRole := c.GetString("user_role")

	var req service.GetAttendancesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	// If user is not admin, only show their own attendances
	if userRole != "admin" {
		req.UserID = &userID
	}

	attendances, total, err := h.service.GetAttendances(c.Request.Context(), req)
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

	utils.PaginatedResponse(c, attendances, page, perPage, total)
}

// CreateManualAttendance handles manual attendance entry
func (h *AttendanceHandler) CreateManualAttendance(c *gin.Context) {
	createdBy := c.GetString("user_id")

	var req service.ManualAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.CreateManualAttendance(c.Request.Context(), createdBy, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 201, result)
}

// UpdateAttendance handles update attendance request
func (h *AttendanceHandler) UpdateAttendance(c *gin.Context) {
	attendanceID := c.Param("id")

	var req service.UpdateAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.UpdateAttendance(c.Request.Context(), attendanceID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GetStatistics handles get attendance statistics request
func (h *AttendanceHandler) GetStatistics(c *gin.Context) {
	userID := c.GetString("user_id")
	userRole := c.GetString("user_role")

	// If admin/dosen, allow querying other users
	queryUserID := c.Query("user_id")
	if queryUserID != "" && (userRole == "dosen" || userRole == "staff") {
		userID = queryUserID
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	var startDatePtr, endDatePtr *string
	if startDate != "" {
		startDatePtr = &startDate
	}
	if endDate != "" {
		endDatePtr = &endDate
	}

	result, err := h.service.GetAttendanceStatistics(c.Request.Context(), userID, startDatePtr, endDatePtr)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GetOverview handles get attendance overview request
func (h *AttendanceHandler) GetOverview(c *gin.Context) {
	userID := c.GetString("user_id")
	userRole := c.GetString("user_role")

	result, err := h.service.GetAttendanceOverview(c.Request.Context(), userID, userRole)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GetByCourse handles get attendance by course request
func (h *AttendanceHandler) GetByCourse(c *gin.Context) {
	courseID := c.Param("courseId")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	var startDatePtr, endDatePtr *string
	if startDate != "" {
		startDatePtr = &startDate
	}
	if endDate != "" {
		endDatePtr = &endDate
	}

	result, err := h.service.GetAttendanceByCourse(c.Request.Context(), courseID, startDatePtr, endDatePtr)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GetByStudent handles get attendance by student request
func (h *AttendanceHandler) GetByStudent(c *gin.Context) {
	studentID := c.Param("studentId")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	var startDatePtr, endDatePtr *string
	if startDate != "" {
		startDatePtr = &startDate
	}
	if endDate != "" {
		endDatePtr = &endDate
	}

	result, err := h.service.GetAttendanceByStudent(c.Request.Context(), studentID, startDatePtr, endDatePtr)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// TapIn handles tap in request for campus attendance
func (h *AttendanceHandler) TapIn(c *gin.Context) {
	userID := c.GetString("user_id")

	var req service.TapInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.TapIn(c.Request.Context(), userID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 201, result)
}

// TapOut handles tap out request for campus attendance
func (h *AttendanceHandler) TapOut(c *gin.Context) {
	userID := c.GetString("user_id")

	var req service.TapInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.TapOut(c.Request.Context(), userID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// CreateSchedule handles create schedule request
func (h *AttendanceHandler) CreateSchedule(c *gin.Context) {
	var req service.CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.CreateSchedule(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 201, result)
}

// GetSchedules handles get schedules request
func (h *AttendanceHandler) GetSchedules(c *gin.Context) {
	var req service.GetSchedulesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	schedules, total, err := h.service.GetSchedules(c.Request.Context(), req)
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

	utils.PaginatedResponse(c, schedules, page, perPage, total)
}

// GetSchedule handles get schedule by ID request
func (h *AttendanceHandler) GetSchedule(c *gin.Context) {
	scheduleID := c.Param("id")

	result, err := h.service.GetScheduleByID(c.Request.Context(), scheduleID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// UpdateSchedule handles update schedule request
func (h *AttendanceHandler) UpdateSchedule(c *gin.Context) {
	scheduleID := c.Param("id")

	var req service.UpdateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.UpdateSchedule(c.Request.Context(), scheduleID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// DeleteSchedule handles delete schedule request
func (h *AttendanceHandler) DeleteSchedule(c *gin.Context) {
	scheduleID := c.Param("id")

	err := h.service.DeleteSchedule(c.Request.Context(), scheduleID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, gin.H{"message": "Schedule deleted successfully"})
}

// GetTodaySchedules handles get today's schedules request
func (h *AttendanceHandler) GetTodaySchedules(c *gin.Context) {
	userID := c.GetString("user_id")
	_ = c.GetString("user_role") // Reserved for future use

	// We need to expose repo through service or add method to service
	// For now, using GetSchedules with today's date filter
	today := time.Now().Format("2006-01-02")
	req := service.GetSchedulesRequest{
		DosenID:   &userID,
		StartDate: &today,
		EndDate:   &today,
	}

	schedules, _, err := h.service.GetSchedules(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, schedules)
}


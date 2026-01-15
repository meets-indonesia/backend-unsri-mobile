package handler

import (
	"net/http"
	"time"

	"unsri-backend/internal/attendance/service"
	"unsri-backend/internal/shared/logger"
	"unsri-backend/internal/shared/utils"

	"github.com/gin-gonic/gin"
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
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.GenerateQRCode(c.Request.Context(), userID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}

// ScanQR handles QR code scan request
func (h *AttendanceHandler) ScanQR(c *gin.Context) {
	userID := c.GetString("user_id")

	var req service.ScanQRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.ScanQRCode(c.Request.Context(), userID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}

// GetAttendances handles get attendances request
func (h *AttendanceHandler) GetAttendances(c *gin.Context) {
	userID := c.GetString("user_id")
	userRole := c.GetString("user_role")

	var req service.GetAttendancesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
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
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.CreateManualAttendance(c.Request.Context(), createdBy, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, result)
}

// UpdateAttendance handles update attendance request
func (h *AttendanceHandler) UpdateAttendance(c *gin.Context) {
	attendanceID := c.Param("id")

	var req service.UpdateAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.UpdateAttendance(c.Request.Context(), attendanceID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
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

	utils.SuccessResponse(c, http.StatusOK, result)
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

	utils.SuccessResponse(c, http.StatusOK, result)
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

	utils.SuccessResponse(c, http.StatusOK, result)
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

	utils.SuccessResponse(c, http.StatusOK, result)
}

// TapIn handles tap in request for campus attendance
func (h *AttendanceHandler) TapIn(c *gin.Context) {
	userID := c.GetString("user_id")

	var req service.TapInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.TapIn(c.Request.Context(), userID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, result)
}

// TapOut handles tap out request for campus attendance
func (h *AttendanceHandler) TapOut(c *gin.Context) {
	userID := c.GetString("user_id")

	var req service.TapInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.TapOut(c.Request.Context(), userID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}

// CreateSchedule handles create schedule request
func (h *AttendanceHandler) CreateSchedule(c *gin.Context) {
	var req service.CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.CreateSchedule(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, result)
}

// GetSchedules handles get schedules request
func (h *AttendanceHandler) GetSchedules(c *gin.Context) {
	var req service.GetSchedulesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
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

	utils.SuccessResponse(c, http.StatusOK, result)
}

// UpdateSchedule handles update schedule request
func (h *AttendanceHandler) UpdateSchedule(c *gin.Context) {
	scheduleID := c.Param("id")

	var req service.UpdateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.UpdateSchedule(c.Request.Context(), scheduleID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
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

// ========== Work Attendance (HRIS) Handlers ==========

// CreateShiftPattern handles create shift pattern request
func (h *AttendanceHandler) CreateShiftPattern(c *gin.Context) {
	var req service.CreateShiftPatternRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.CreateShiftPattern(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, result)
}

// GetShiftPattern handles get shift pattern by ID request
func (h *AttendanceHandler) GetShiftPattern(c *gin.Context) {
	shiftID := c.Param("id")

	result, err := h.service.GetShiftPatternByID(c.Request.Context(), shiftID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}

// GetShiftPatterns handles get shift patterns request
func (h *AttendanceHandler) GetShiftPatterns(c *gin.Context) {
	var req service.GetShiftPatternsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	shifts, total, err := h.service.GetShiftPatterns(c.Request.Context(), req)
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

	utils.PaginatedResponse(c, shifts, page, perPage, total)
}

// UpdateShiftPattern handles update shift pattern request
func (h *AttendanceHandler) UpdateShiftPattern(c *gin.Context) {
	shiftID := c.Param("id")

	var req service.UpdateShiftPatternRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.UpdateShiftPattern(c.Request.Context(), shiftID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}

// DeleteShiftPattern handles delete shift pattern request
func (h *AttendanceHandler) DeleteShiftPattern(c *gin.Context) {
	shiftID := c.Param("id")

	err := h.service.DeleteShiftPattern(c.Request.Context(), shiftID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, gin.H{"message": "Shift pattern deleted successfully"})
}

// CreateUserShift handles create user shift request
func (h *AttendanceHandler) CreateUserShift(c *gin.Context) {
	var req service.CreateUserShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.CreateUserShift(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, result)
}

// GetUserShifts handles get user shifts request
func (h *AttendanceHandler) GetUserShifts(c *gin.Context) {
	userID := c.Param("userId")
	date := c.Query("date")

	var datePtr *string
	if date != "" {
		datePtr = &date
	}

	result, err := h.service.GetUserShiftsByUserID(c.Request.Context(), userID, datePtr)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}

// CreateWorkSchedule handles create work schedule request
func (h *AttendanceHandler) CreateWorkSchedule(c *gin.Context) {
	var req service.CreateWorkScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.CreateWorkSchedule(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, result)
}

// GetWorkSchedules handles get work schedules request
func (h *AttendanceHandler) GetWorkSchedules(c *gin.Context) {
	var req service.GetWorkSchedulesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	schedules, total, err := h.service.GetWorkSchedules(c.Request.Context(), req)
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

// CheckIn handles work check-in request
func (h *AttendanceHandler) CheckIn(c *gin.Context) {
	userID := c.GetString("user_id")

	var req service.CheckInRequest
	if err := c.ShouldBind(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.CheckIn(c.Request.Context(), userID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, result)
}

// CheckOut handles work check-out request
func (h *AttendanceHandler) CheckOut(c *gin.Context) {
	userID := c.GetString("user_id")

	var req service.CheckOutRequest
	if err := c.ShouldBind(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.CheckOut(c.Request.Context(), userID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, result)
}

// GetWorkAttendanceRecords handles get work attendance records request
func (h *AttendanceHandler) GetWorkAttendanceRecords(c *gin.Context) {
	var req service.GetWorkAttendanceRecordsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// If user_id not provided, use current user
	if req.UserID == "" {
		req.UserID = c.GetString("user_id")
	}

	records, total, err := h.service.GetWorkAttendanceRecords(c.Request.Context(), req)
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

	utils.PaginatedResponse(c, records, page, perPage, total)
}

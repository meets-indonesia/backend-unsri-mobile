package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"unsri-backend/internal/schedule/service"
	"unsri-backend/internal/shared/logger"
	"unsri-backend/internal/shared/utils"
)

// ScheduleHandler handles HTTP requests for schedule
type ScheduleHandler struct {
	service *service.ScheduleService
	logger  logger.Logger
}

// NewScheduleHandler creates a new schedule handler
func NewScheduleHandler(service *service.ScheduleService, logger logger.Logger) *ScheduleHandler {
	return &ScheduleHandler{
		service: service,
		logger:  logger,
	}
}

// CreateSchedule handles create schedule request
func (h *ScheduleHandler) CreateSchedule(c *gin.Context) {
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

// GetSchedule handles get schedule by ID request
func (h *ScheduleHandler) GetSchedule(c *gin.Context) {
	id := c.Param("id")

	result, err := h.service.GetScheduleByID(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GetSchedules handles get schedules request
func (h *ScheduleHandler) GetSchedules(c *gin.Context) {
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

// GetTodaySchedules handles get today's schedules request
func (h *ScheduleHandler) GetTodaySchedules(c *gin.Context) {
	userID := c.GetString("user_id")
	userRole := c.GetString("user_role")

	result, err := h.service.GetTodaySchedules(c.Request.Context(), userID, userRole)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GetUpcomingSchedules handles get upcoming schedules request
func (h *ScheduleHandler) GetUpcomingSchedules(c *gin.Context) {
	userID := c.GetString("user_id")
	userRole := c.GetString("user_role")
	limit := 10

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	result, err := h.service.GetUpcomingSchedules(c.Request.Context(), userID, userRole, limit)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GetCalendarView handles get calendar view request
func (h *ScheduleHandler) GetCalendarView(c *gin.Context) {
	userID := c.GetString("user_id")
	userRole := c.GetString("user_role")

	year, _ := strconv.Atoi(c.Param("year"))
	month, _ := strconv.Atoi(c.Param("month"))

	if year == 0 {
		year = 2024
	}
	if month == 0 || month < 1 || month > 12 {
		month = 1
	}

	result, err := h.service.GetCalendarView(c.Request.Context(), userID, userRole, year, month)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// UpdateSchedule handles update schedule request
func (h *ScheduleHandler) UpdateSchedule(c *gin.Context) {
	id := c.Param("id")

	var req service.UpdateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.UpdateSchedule(c.Request.Context(), id, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// DeleteSchedule handles delete schedule request
func (h *ScheduleHandler) DeleteSchedule(c *gin.Context) {
	id := c.Param("id")

	err := h.service.DeleteSchedule(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, gin.H{"message": "Schedule deleted successfully"})
}


package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"unsri-backend/internal/calendar/service"
	"unsri-backend/internal/shared/logger"
	"unsri-backend/internal/shared/utils"
)

// CalendarHandler handles HTTP requests for calendar
type CalendarHandler struct {
	service *service.CalendarService
	logger  logger.Logger
}

// NewCalendarHandler creates a new calendar handler
func NewCalendarHandler(service *service.CalendarService, logger logger.Logger) *CalendarHandler {
	return &CalendarHandler{
		service: service,
		logger:  logger,
	}
}

// CreateEvent handles create event request
func (h *CalendarHandler) CreateEvent(c *gin.Context) {
	createdBy := c.GetString("user_id")

	var req service.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.CreateEvent(c.Request.Context(), createdBy, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 201, result)
}

// GetEvent handles get event by ID request
func (h *CalendarHandler) GetEvent(c *gin.Context) {
	id := c.Param("id")

	result, err := h.service.GetEventByID(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GetEvents handles get events request
func (h *CalendarHandler) GetEvents(c *gin.Context) {
	var req service.GetEventsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	events, total, err := h.service.GetEvents(c.Request.Context(), req)
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

	utils.PaginatedResponse(c, events, page, perPage, total)
}

// GetEventsByMonth handles get events by month request
func (h *CalendarHandler) GetEventsByMonth(c *gin.Context) {
	year, _ := strconv.Atoi(c.Param("year"))
	month, _ := strconv.Atoi(c.Param("month"))

	if year == 0 {
		year = 2024
	}
	if month == 0 || month < 1 || month > 12 {
		month = 1
	}

	result, err := h.service.GetEventsByMonth(c.Request.Context(), year, month)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GetUpcomingEvents handles get upcoming events request
func (h *CalendarHandler) GetUpcomingEvents(c *gin.Context) {
	limit := 10

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	result, err := h.service.GetUpcomingEvents(c.Request.Context(), limit)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// UpdateEvent handles update event request
func (h *CalendarHandler) UpdateEvent(c *gin.Context) {
	id := c.Param("id")

	var req service.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.UpdateEvent(c.Request.Context(), id, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// DeleteEvent handles delete event request
func (h *CalendarHandler) DeleteEvent(c *gin.Context) {
	id := c.Param("id")

	err := h.service.DeleteEvent(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, gin.H{"message": "Event deleted successfully"})
}


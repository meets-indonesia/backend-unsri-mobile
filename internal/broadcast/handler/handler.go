package handler

import (
	"net/http"
	"strconv"

	"unsri-backend/internal/broadcast/service"
	"unsri-backend/internal/shared/logger"
	"unsri-backend/internal/shared/utils"

	"github.com/gin-gonic/gin"
)

// BroadcastHandler handles HTTP requests for broadcast
type BroadcastHandler struct {
	service *service.BroadcastService
	logger  logger.Logger
}

// NewBroadcastHandler creates a new broadcast handler
func NewBroadcastHandler(service *service.BroadcastService, logger logger.Logger) *BroadcastHandler {
	return &BroadcastHandler{
		service: service,
		logger:  logger,
	}
}

// CreateBroadcast handles create broadcast request
func (h *BroadcastHandler) CreateBroadcast(c *gin.Context) {
	createdBy := c.GetString("user_id")

	var req service.CreateBroadcastRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.CreateBroadcast(c.Request.Context(), createdBy, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, result)
}

// GetBroadcast handles get broadcast by ID request
func (h *BroadcastHandler) GetBroadcast(c *gin.Context) {
	id := c.Param("id")

	result, err := h.service.GetBroadcastByID(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}

// GetBroadcasts handles get broadcasts request
func (h *BroadcastHandler) GetBroadcasts(c *gin.Context) {
	var req service.GetBroadcastsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	broadcasts, total, err := h.service.GetBroadcasts(c.Request.Context(), req)
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

	utils.PaginatedResponse(c, broadcasts, page, perPage, total)
}

// UpdateBroadcast handles update broadcast request
func (h *BroadcastHandler) UpdateBroadcast(c *gin.Context) {
	id := c.Param("id")

	var req service.UpdateBroadcastRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.UpdateBroadcast(c.Request.Context(), id, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}

// DeleteBroadcast handles delete broadcast request
func (h *BroadcastHandler) DeleteBroadcast(c *gin.Context) {
	id := c.Param("id")

	err := h.service.DeleteBroadcast(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Broadcast deleted successfully"})
}

// SearchBroadcasts handles search broadcasts request
func (h *BroadcastHandler) SearchBroadcasts(c *gin.Context) {
	query := c.Query("q")
	page := 1
	perPage := 20

	if p := c.Query("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}
	if pp := c.Query("per_page"); pp != "" {
		if val, err := strconv.Atoi(pp); err == nil && val > 0 {
			perPage = val
		}
	}

	result, total, err := h.service.SearchBroadcasts(c.Request.Context(), query, page, perPage)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.PaginatedResponse(c, result, page, perPage, total)
}

// GetGeneralBroadcasts handles get general broadcasts request
func (h *BroadcastHandler) GetGeneralBroadcasts(c *gin.Context) {
	page := 1
	perPage := 20

	result, total, err := h.service.GetGeneralBroadcasts(c.Request.Context(), page, perPage)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.PaginatedResponse(c, result, page, perPage, total)
}

// GetClassBroadcasts handles get class broadcasts request
func (h *BroadcastHandler) GetClassBroadcasts(c *gin.Context) {
	classID := c.Query("class_id")
	page := 1
	perPage := 20

	var classIDPtr *string
	if classID != "" {
		classIDPtr = &classID
	}

	result, total, err := h.service.GetClassBroadcasts(c.Request.Context(), classIDPtr, page, perPage)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.PaginatedResponse(c, result, page, perPage, total)
}

// ScheduleBroadcast handles schedule broadcast request
func (h *BroadcastHandler) ScheduleBroadcast(c *gin.Context) {
	id := c.Param("id")

	var req service.ScheduleBroadcastRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.ScheduleBroadcast(c.Request.Context(), id, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, result)
}

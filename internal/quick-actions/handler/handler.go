package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"unsri-backend/internal/quick-actions/service"
	"unsri-backend/internal/shared/logger"
	"unsri-backend/internal/shared/utils"
)

// QuickActionsHandler handles HTTP requests for quick actions
type QuickActionsHandler struct {
	service *service.QuickActionsService
	logger  logger.Logger
}

// NewQuickActionsHandler creates a new quick actions handler
func NewQuickActionsHandler(service *service.QuickActionsService, logger logger.Logger) *QuickActionsHandler {
	return &QuickActionsHandler{
		service: service,
		logger:  logger,
	}
}

// GetQuickActions handles get quick actions request
func (h *QuickActionsHandler) GetQuickActions(c *gin.Context) {
	userRole := c.GetString("user_role")

	result, err := h.service.GetQuickActions(c.Request.Context(), userRole)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GetTranscript handles get transcript request
func (h *QuickActionsHandler) GetTranscript(c *gin.Context) {
	studentID := c.Param("studentId")

	result, err := h.service.GetTranscript(c.Request.Context(), studentID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GetKRS handles get KRS request
func (h *QuickActionsHandler) GetKRS(c *gin.Context) {
	studentID := c.Param("studentId")
	semester := c.Query("semester")

	var semesterPtr *string
	if semester != "" {
		semesterPtr = &semester
	}

	result, err := h.service.GetKRS(c.Request.Context(), studentID, semesterPtr)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GetBimbingans handles get bimbingans request
func (h *QuickActionsHandler) GetBimbingans(c *gin.Context) {
	userID := c.GetString("user_id")
	userRole := c.GetString("user_role")
	limit := 10

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	result, err := h.service.GetBimbingans(c.Request.Context(), userID, userRole, limit)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}


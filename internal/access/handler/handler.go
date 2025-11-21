package handler

import (
	"github.com/gin-gonic/gin"
	"unsri-backend/internal/access/service"
	"unsri-backend/internal/shared/logger"
	"unsri-backend/internal/shared/utils"
)

// AccessHandler handles HTTP requests for access control
type AccessHandler struct {
	service *service.AccessService
	logger  logger.Logger
}

// NewAccessHandler creates a new access handler
func NewAccessHandler(service *service.AccessService, logger logger.Logger) *AccessHandler {
	return &AccessHandler{
		service: service,
		logger:  logger,
	}
}

// ValidateQR handles validate access QR request
func (h *AccessHandler) ValidateQR(c *gin.Context) {
	var req service.ValidateQRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.ValidateAccessQR(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GetAccessHistory handles get access history request
func (h *AccessHandler) GetAccessHistory(c *gin.Context) {
	var req service.GetAccessHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	logs, total, err := h.service.GetAccessHistory(c.Request.Context(), req)
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

	utils.PaginatedResponse(c, logs, page, perPage, total)
}

// LogAccess handles log access request
func (h *AccessHandler) LogAccess(c *gin.Context) {
	var req service.LogAccessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.LogAccess(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 201, result)
}

// GetAccessPermissions handles get access permissions request
func (h *AccessHandler) GetAccessPermissions(c *gin.Context) {
	userID := c.Param("userId")

	result, err := h.service.GetAccessPermissions(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// CreateAccessPermission handles create access permission request
func (h *AccessHandler) CreateAccessPermission(c *gin.Context) {
	var req service.CreateAccessPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.CreateAccessPermission(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 201, result)
}


package handler

import (
	"github.com/gin-gonic/gin"
	"unsri-backend/internal/qr/service"
	"unsri-backend/internal/shared/logger"
	"unsri-backend/internal/shared/utils"
)

// QRHandler handles HTTP requests for QR service
type QRHandler struct {
	service *service.QRService
	logger  logger.Logger
}

// NewQRHandler creates a new QR handler
func NewQRHandler(service *service.QRService, logger logger.Logger) *QRHandler {
	return &QRHandler{
		service: service,
		logger:  logger,
	}
}

// GenerateQR handles generate QR request
func (h *QRHandler) GenerateQR(c *gin.Context) {
	createdBy := c.GetString("user_id")

	var req service.GenerateQRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.GenerateQR(c.Request.Context(), createdBy, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 201, result)
}

// ValidateQR handles validate QR request
func (h *QRHandler) ValidateQR(c *gin.Context) {
	var req service.ValidateQRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.ValidateQR(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GetQR handles get QR by ID request
func (h *QRHandler) GetQR(c *gin.Context) {
	id := c.Param("id")

	result, err := h.service.GetQRByID(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GenerateClassQR handles generate class QR request
func (h *QRHandler) GenerateClassQR(c *gin.Context) {
	createdBy := c.GetString("user_id")

	var req service.GenerateClassQRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.GenerateClassQR(c.Request.Context(), createdBy, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 201, result)
}

// GenerateAccessQR handles generate access QR request (gate access - unique per user)
func (h *QRHandler) GenerateAccessQR(c *gin.Context) {
	userID := c.GetString("user_id")

	result, err := h.service.GenerateAccessQR(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// ValidateAccessQR handles validate access QR request (for gate)
func (h *QRHandler) ValidateAccessQR(c *gin.Context) {
	qrToken := c.Param("token")

	result, err := h.service.ValidateAccessQR(c.Request.Context(), qrToken)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// RegenerateClassQR handles regenerate class QR request (after scan)
func (h *QRHandler) RegenerateClassQR(c *gin.Context) {
	scheduleID := c.Param("scheduleId")
	createdBy := c.GetString("user_id")

	result, err := h.service.RegenerateClassQR(c.Request.Context(), scheduleID, createdBy)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}


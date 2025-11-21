package handler

import (
	"github.com/gin-gonic/gin"
	"unsri-backend/internal/notification/service"
	"unsri-backend/internal/shared/logger"
	"unsri-backend/internal/shared/utils"
)

// NotificationHandler handles HTTP requests for notification
type NotificationHandler struct {
	service *service.NotificationService
	logger  logger.Logger
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(service *service.NotificationService, logger logger.Logger) *NotificationHandler {
	return &NotificationHandler{
		service: service,
		logger:  logger,
	}
}

// SendNotification handles send notification request
func (h *NotificationHandler) SendNotification(c *gin.Context) {
	var req service.SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.SendNotification(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 201, result)
}

// GetNotifications handles get notifications request
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID := c.GetString("user_id")

	var req service.GetNotificationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	notifications, total, err := h.service.GetNotifications(c.Request.Context(), userID, req)
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

	utils.PaginatedResponse(c, notifications, page, perPage, total)
}

// MarkAsRead handles mark notification as read request
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	notificationID := c.Param("id")
	userID := c.GetString("user_id")

	err := h.service.MarkAsRead(c.Request.Context(), notificationID, userID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, gin.H{"message": "Notification marked as read"})
}

// MarkAllAsRead handles mark all notifications as read request
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := c.GetString("user_id")

	err := h.service.MarkAllAsRead(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, gin.H{"message": "All notifications marked as read"})
}

// RegisterDeviceToken handles register device token request
func (h *NotificationHandler) RegisterDeviceToken(c *gin.Context) {
	userID := c.GetString("user_id")

	var req service.RegisterDeviceTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.RegisterDeviceToken(c.Request.Context(), userID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 201, result)
}

// UnregisterDeviceToken handles unregister device token request
func (h *NotificationHandler) UnregisterDeviceToken(c *gin.Context) {
	token := c.Param("token")

	err := h.service.UnregisterDeviceToken(c.Request.Context(), token)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, gin.H{"message": "Device token unregistered"})
}


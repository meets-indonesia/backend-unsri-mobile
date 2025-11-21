package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"unsri-backend/internal/auth/service"
	"unsri-backend/internal/shared/logger"
	"unsri-backend/internal/shared/utils"
)

// AuthHandler handles HTTP requests for authentication
type AuthHandler struct {
	service *service.AuthService
	logger  logger.Logger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(service *service.AuthService, logger logger.Logger) *AuthHandler {
	return &AuthHandler{
		service: service,
		logger:  logger,
	}
}

// Login handles login request
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// Register handles registration request
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 201, result)
}

// RefreshToken handles refresh token request
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req service.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.RefreshToken(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// VerifyToken handles token verification request
func (h *AuthHandler) VerifyToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		utils.ErrorResponse(c, 401, errors.New("authorization header required"))
		return
	}

	// Remove "Bearer " prefix if present
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	result, err := h.service.VerifyToken(c.Request.Context(), token)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}


package handler

import (
	"github.com/gin-gonic/gin"
	"unsri-backend/internal/shared/logger"
	"unsri-backend/internal/shared/utils"
	"unsri-backend/internal/user/service"
)

// UserHandler handles HTTP requests for user management
type UserHandler struct {
	service *service.UserService
	logger  logger.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(service *service.UserService, logger logger.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

// GetProfile handles get user profile request
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	result, err := h.service.GetUserProfile(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// UpdateProfile handles update user profile request
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	var req service.UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	result, err := h.service.UpdateUserProfile(c.Request.Context(), userID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GetUserByID handles get user by ID request
func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")

	result, err := h.service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// SearchUsers handles search users request
func (h *UserHandler) SearchUsers(c *gin.Context) {
	var req service.SearchUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	users, total, err := h.service.SearchUsers(c.Request.Context(), req)
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

	utils.PaginatedResponse(c, users, page, perPage, total)
}

// GetMahasiswaByNIM handles get mahasiswa by NIM request
func (h *UserHandler) GetMahasiswaByNIM(c *gin.Context) {
	nim := c.Param("nim")

	result, err := h.service.GetMahasiswaByNIM(c.Request.Context(), nim)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GetDosenByNIP handles get dosen by NIP request
func (h *UserHandler) GetDosenByNIP(c *gin.Context) {
	nip := c.Param("nip")

	result, err := h.service.GetDosenByNIP(c.Request.Context(), nip)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GetStaffByNIP handles get staff by NIP request
func (h *UserHandler) GetStaffByNIP(c *gin.Context) {
	nip := c.Param("nip")

	result, err := h.service.GetStaffByNIP(c.Request.Context(), nip)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// UploadAvatar handles avatar upload request
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	userID := c.GetString("user_id")

	file, err := c.FormFile("avatar")
	if err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	// Read file
	src, err := file.Open()
	if err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}
	defer src.Close()

	data := make([]byte, file.Size)
	if _, err := src.Read(data); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	req := service.UploadAvatarRequest{
		Filename: file.Filename,
		Data:     data,
		MimeType: file.Header.Get("Content-Type"),
	}

	avatarURL, err := h.service.UploadAvatar(c.Request.Context(), userID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, gin.H{"avatar_url": avatarURL})
}


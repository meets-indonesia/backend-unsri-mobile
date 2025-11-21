package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"unsri-backend/internal/file-storage/service"
	"unsri-backend/internal/shared/logger"
	"unsri-backend/internal/shared/utils"
)

// FileStorageHandler handles HTTP requests for file storage
type FileStorageHandler struct {
	service *service.FileStorageService
	logger  logger.Logger
}

// NewFileStorageHandler creates a new file storage handler
func NewFileStorageHandler(service *service.FileStorageService, logger logger.Logger) *FileStorageHandler {
	return &FileStorageHandler{
		service: service,
		logger:  logger,
	}
}

// UploadFile handles file upload request
func (h *FileStorageHandler) UploadFile(c *gin.Context) {
	userID := c.GetString("user_id")
	fileType := c.PostForm("file_type")
	if fileType == "" {
		fileType = "document"
	}

	file, err := c.FormFile("file")
	if err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	req := service.UploadFileRequest{
		File:     file,
		FileType: fileType,
		IsPublic: c.PostForm("is_public") == "true",
	}

	result, err := h.service.UploadFile(c.Request.Context(), userID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 201, result)
}

// GetFile handles get file request
func (h *FileStorageHandler) GetFile(c *gin.Context) {
	id := c.Param("id")

	result, err := h.service.GetFileByID(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GetFiles handles get files request
func (h *FileStorageHandler) GetFiles(c *gin.Context) {
	userID := c.GetString("user_id")

	var req service.GetFilesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	files, total, err := h.service.GetFiles(c.Request.Context(), userID, req)
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

	utils.PaginatedResponse(c, files, page, perPage, total)
}

// DeleteFile handles delete file request
func (h *FileStorageHandler) DeleteFile(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")

	err := h.service.DeleteFile(c.Request.Context(), id, userID)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, gin.H{"message": "File deleted successfully"})
}

// DownloadFile handles file download request
func (h *FileStorageHandler) DownloadFile(c *gin.Context) {
	id := c.Param("id")

	content, mimeType, err := h.service.GetFileContent(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	c.Data(http.StatusOK, mimeType, content)
}

// UploadAvatar handles avatar upload request
func (h *FileStorageHandler) UploadAvatar(c *gin.Context) {
	userID := c.GetString("user_id")

	file, err := c.FormFile("file")
	if err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	req := service.UploadAvatarRequest{
		File: file,
	}

	result, err := h.service.UploadAvatar(c.Request.Context(), userID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 201, result)
}

// UploadDocument handles document upload request
func (h *FileStorageHandler) UploadDocument(c *gin.Context) {
	userID := c.GetString("user_id")

	file, err := c.FormFile("file")
	if err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	req := service.UploadDocumentRequest{
		File:     file,
		IsPublic: c.PostForm("is_public") == "true",
	}

	result, err := h.service.UploadDocument(c.Request.Context(), userID, req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 201, result)
}


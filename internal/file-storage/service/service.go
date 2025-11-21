package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	apperrors "unsri-backend/internal/shared/errors"
	"unsri-backend/internal/shared/models"
	"unsri-backend/internal/file-storage/repository"
	"github.com/google/uuid"
)

// FileStorageService handles file storage business logic
type FileStorageService struct {
	repo   *repository.FileRepository
	config StorageConfig
}

// StorageConfig holds storage configuration
type StorageConfig struct {
	Type     string
	BasePath string
	BaseURL  string
	MaxSize  int64
}

// NewFileStorageService creates a new file storage service
func NewFileStorageService(repo *repository.FileRepository, config StorageConfig) *FileStorageService {
	return &FileStorageService{
		repo:   repo,
		config: config,
	}
}

// UploadFileRequest represents upload file request
type UploadFileRequest struct {
	File     *multipart.FileHeader
	FileType string
	IsPublic bool
}

// UploadFile uploads a file
func (s *FileStorageService) UploadFile(ctx context.Context, userID string, req UploadFileRequest) (*models.File, error) {
	// Validate file size
	if req.File.Size > s.config.MaxSize {
		return nil, apperrors.NewValidationError(fmt.Sprintf("file size exceeds maximum allowed size of %d bytes", s.config.MaxSize))
	}

	// Generate unique filename
	ext := filepath.Ext(req.File.Filename)
	fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	
	// Create directory if not exists
	fileDir := filepath.Join(s.config.BasePath, req.FileType)
	if err := os.MkdirAll(fileDir, 0755); err != nil {
		return nil, apperrors.NewInternalError("failed to create directory", err)
	}

	// Save file
	filePath := filepath.Join(fileDir, fileName)
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to create file", err)
	}
	defer dst.Close()

	src, err := req.File.Open()
	if err != nil {
		return nil, apperrors.NewInternalError("failed to open file", err)
	}
	defer src.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, apperrors.NewInternalError("failed to save file", err)
	}

	// Create file record
	file := &models.File{
		UserID:       userID,
		FileName:     fileName,
		OriginalName: req.File.Filename,
		FileType:     req.FileType,
		MimeType:     req.File.Header.Get("Content-Type"),
		Size:         req.File.Size,
		Path:         filePath,
		URL:          fmt.Sprintf("%s/%s/%s", s.config.BaseURL, req.FileType, fileName),
		IsPublic:     req.IsPublic,
	}

	if err := s.repo.CreateFile(ctx, file); err != nil {
		os.Remove(filePath) // Cleanup on error
		return nil, apperrors.NewInternalError("failed to create file record", err)
	}

	return file, nil
}

// GetFileByID gets a file by ID
func (s *FileStorageService) GetFileByID(ctx context.Context, id string) (*models.File, error) {
	return s.repo.GetFileByID(ctx, id)
}

// GetFilesRequest represents get files request
type GetFilesRequest struct {
	FileType string `form:"file_type"`
	Page     int    `form:"page,default=1"`
	PerPage  int    `form:"per_page,default=20"`
}

// GetFiles gets files for a user
func (s *FileStorageService) GetFiles(ctx context.Context, userID string, req GetFilesRequest) ([]models.File, int64, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}

	var fileTypePtr *string
	if req.FileType != "" {
		fileTypePtr = &req.FileType
	}

	return s.repo.GetFilesByUserID(ctx, userID, fileTypePtr, perPage, (page-1)*perPage)
}

// DeleteFile deletes a file
func (s *FileStorageService) DeleteFile(ctx context.Context, id string, userID string) error {
	file, err := s.repo.GetFileByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFoundError("file", id)
	}

	// Check ownership
	if file.UserID != userID {
		return apperrors.NewForbiddenError("not authorized to delete this file")
	}

	// Delete physical file
	if err := os.Remove(file.Path); err != nil && !os.IsNotExist(err) {
		return apperrors.NewInternalError("failed to delete file", err)
	}

	// Delete record
	return s.repo.DeleteFile(ctx, id)
}

// UploadAvatarRequest represents upload avatar request
type UploadAvatarRequest struct {
	File *multipart.FileHeader
}

// UploadAvatar uploads user avatar
func (s *FileStorageService) UploadAvatar(ctx context.Context, userID string, req UploadAvatarRequest) (*models.File, error) {
	// Delete old avatar
	oldAvatar, _ := s.repo.GetAvatarByUserID(ctx, userID)
	if oldAvatar != nil {
		os.Remove(oldAvatar.Path)
		s.repo.DeleteFile(ctx, oldAvatar.ID)
	}

	uploadReq := UploadFileRequest{
		File:     req.File,
		FileType: "avatar",
		IsPublic: true,
	}

	return s.UploadFile(ctx, userID, uploadReq)
}

// UploadDocumentRequest represents upload document request
type UploadDocumentRequest struct {
	File     *multipart.FileHeader
	IsPublic bool
}

// UploadDocument uploads a document
func (s *FileStorageService) UploadDocument(ctx context.Context, userID string, req UploadDocumentRequest) (*models.File, error) {
	uploadReq := UploadFileRequest{
		File:     req.File,
		FileType: "document",
		IsPublic: req.IsPublic,
	}

	return s.UploadFile(ctx, userID, uploadReq)
}

// GetFileContent gets file content for download
func (s *FileStorageService) GetFileContent(ctx context.Context, id string) ([]byte, string, error) {
	file, err := s.repo.GetFileByID(ctx, id)
	if err != nil {
		return nil, "", apperrors.NewNotFoundError("file", id)
	}

	content, err := os.ReadFile(file.Path)
	if err != nil {
		return nil, "", apperrors.NewInternalError("failed to read file", err)
	}

	return content, file.MimeType, nil
}


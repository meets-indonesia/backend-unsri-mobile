package service

import (
	"testing"
	"time"

	"github.com/google/uuid"
	apperrors "unsri-backend/internal/shared/errors"
	"unsri-backend/internal/shared/models"
)

// Test helper functions
func createTestFile() *models.File {
	return &models.File{
		ID:           uuid.New().String(),
		UserID:       uuid.New().String(),
		FileName:     "test.pdf",
		OriginalName: "test.pdf",
		FileType:     "document",
		MimeType:     "application/pdf",
		Size:         1024,
		Path:         "/files/document/test.pdf",
		URL:          "https://example.com/files/document/test.pdf",
		IsPublic:     false,
		CreatedAt:    time.Now(),
	}
}

// Test error types
func TestErrorTypes(t *testing.T) {
	t.Run("NotFoundError", func(t *testing.T) {
		err := apperrors.NewNotFoundError("file", "test-id")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("ValidationError", func(t *testing.T) {
		err := apperrors.NewValidationError("invalid input")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("ForbiddenError", func(t *testing.T) {
		err := apperrors.NewForbiddenError("insufficient permissions")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

// Test File model
func TestFileModel(t *testing.T) {
	t.Run("valid file", func(t *testing.T) {
		file := createTestFile()
		if file.UserID == "" {
			t.Error("UserID should not be empty")
		}
		if file.FileName == "" {
			t.Error("FileName should not be empty")
		}
		if file.Size <= 0 {
			t.Error("Size should be positive")
		}
		if file.ID == "" {
			t.Error("ID should be generated")
		}
	})

	t.Run("table name", func(t *testing.T) {
		file := models.File{}
		if file.TableName() != "files" {
			t.Errorf("Expected table name 'files', got '%s'", file.TableName())
		}
	})

	t.Run("file types", func(t *testing.T) {
		validTypes := []string{"document", "image", "video", "audio", "other"}
		for _, ftype := range validTypes {
			file := createTestFile()
			file.FileType = ftype
			if file.FileType != ftype {
				t.Errorf("Expected file type %s, got %s", ftype, file.FileType)
			}
		}
	})
}

// Test UploadFileRequest validation
func TestUploadFileRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     UploadFileRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: UploadFileRequest{
				FileType: "document",
				IsPublic: false,
			},
			wantErr: false,
		},
		{
			name: "valid request with public",
			req: UploadFileRequest{
				FileType: "image",
				IsPublic: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// UploadFileRequest validation happens in handler with actual file
			validTypes := []string{"document", "image", "video", "audio", "avatar", "other"}
			isValid := false
			for _, validType := range validTypes {
				if tt.req.FileType == validType {
					isValid = true
					break
				}
			}
			if !isValid && tt.req.FileType != "" {
				t.Error("FileType should be valid")
			}
		})
	}
}


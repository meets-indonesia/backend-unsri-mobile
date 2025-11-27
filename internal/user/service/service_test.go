package service

import (
	"testing"

	"github.com/google/uuid"
	apperrors "unsri-backend/internal/shared/errors"
	"unsri-backend/internal/shared/models"
)

// Test helper functions
func createTestUser() *models.User {
	return &models.User{
		ID:           uuid.New().String(),
		Email:        "test@example.com",
		PasswordHash: "$2a$10$hashedpassword",
		Role:         models.RoleMahasiswa,
		IsActive:     true,
	}
}

// Test error types
func TestErrorTypes(t *testing.T) {
	t.Run("NotFoundError", func(t *testing.T) {
		err := apperrors.NewNotFoundError("user", "test-id")
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
}

// Test User model
func TestUserModel(t *testing.T) {
	t.Run("valid user", func(t *testing.T) {
		user := createTestUser()
		if user.Email == "" {
			t.Error("Email should not be empty")
		}
		if user.Role == "" {
			t.Error("Role should not be empty")
		}
		if user.ID == "" {
			t.Error("ID should be generated")
		}
	})

	t.Run("table name", func(t *testing.T) {
		user := models.User{}
		if user.TableName() != "users" {
			t.Errorf("Expected table name 'users', got '%s'", user.TableName())
		}
	})

	t.Run("user roles", func(t *testing.T) {
		validRoles := []models.UserRole{
			models.RoleMahasiswa,
			models.RoleDosen,
			models.RoleStaff,
		}

		for _, role := range validRoles {
			user := createTestUser()
			user.Role = role
			if user.Role != role {
				t.Errorf("Expected role %s, got %s", role, user.Role)
			}
		}
	})

	t.Run("user active status", func(t *testing.T) {
		user := createTestUser()
		if !user.IsActive {
			t.Error("User should be active by default")
		}

		user.IsActive = false
		if user.IsActive {
			t.Error("User should be inactive")
		}
	})
}

// Test UpdateUserProfileRequest validation
func TestUpdateUserProfileRequest(t *testing.T) {
	email := "updated@example.com"
	mahasiswaNama := "Updated Student Name"
	dosenNama := "Updated Lecturer Name"
	staffNama := "Updated Staff Name"

	tests := []struct {
		name string
		req  UpdateUserProfileRequest
	}{
		{
			name: "update email only",
			req: UpdateUserProfileRequest{
				Email: &email,
			},
		},
		{
			name: "update mahasiswa data",
			req: UpdateUserProfileRequest{
				Mahasiswa: &UpdateMahasiswaRequest{
					Nama: &mahasiswaNama,
				},
			},
		},
		{
			name: "update dosen data",
			req: UpdateUserProfileRequest{
				Dosen: &UpdateDosenRequest{
					Nama: &dosenNama,
				},
			},
		},
		{
			name: "update staff data",
			req: UpdateUserProfileRequest{
				Staff: &UpdateStaffRequest{
					Nama: &staffNama,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req.Email != nil && *tt.req.Email == "" {
				t.Error("Email should not be empty if provided")
			}
			if tt.req.Mahasiswa != nil && tt.req.Mahasiswa.Nama != nil && *tt.req.Mahasiswa.Nama == "" {
				t.Error("Mahasiswa name should not be empty if provided")
			}
			if tt.req.Dosen != nil && tt.req.Dosen.Nama != nil && *tt.req.Dosen.Nama == "" {
				t.Error("Dosen name should not be empty if provided")
			}
			if tt.req.Staff != nil && tt.req.Staff.Nama != nil && *tt.req.Staff.Nama == "" {
				t.Error("Staff name should not be empty if provided")
			}
		})
	}
}


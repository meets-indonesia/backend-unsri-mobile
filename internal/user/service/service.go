package service

import (
	"context"
	"errors"
	"strings"

	apperrors "unsri-backend/internal/shared/errors"
	"unsri-backend/internal/shared/models"
	"unsri-backend/internal/user/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService handles user business logic
type UserService struct {
	repo *repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// GetUserProfile gets user profile by ID
func (s *UserService) GetUserProfile(ctx context.Context, userID string) (*models.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

// UpdateUserProfileRequest represents update profile request
type UpdateUserProfileRequest struct {
	Email     *string                 `json:"email,omitempty"`
	Mahasiswa *UpdateMahasiswaRequest `json:"mahasiswa,omitempty"`
	Dosen     *UpdateDosenRequest     `json:"dosen,omitempty"`
	Staff     *UpdateStaffRequest     `json:"staff,omitempty"`
}

// UpdateMahasiswaRequest represents update mahasiswa request
type UpdateMahasiswaRequest struct {
	Nama     *string `json:"nama,omitempty"`
	Prodi    *string `json:"prodi,omitempty"`
	Angkatan *int    `json:"angkatan,omitempty"`
}

// UpdateDosenRequest represents update dosen request
type UpdateDosenRequest struct {
	Nama  *string `json:"nama,omitempty"`
	Prodi *string `json:"prodi,omitempty"`
}

// UpdateStaffRequest represents update staff request
type UpdateStaffRequest struct {
	Nama    *string `json:"nama,omitempty"`
	Jabatan *string `json:"jabatan,omitempty"`
	Unit    *string `json:"unit,omitempty"`
}

// UpdateUserProfile updates user profile
func (s *UserService) UpdateUserProfile(ctx context.Context, userID string, req UpdateUserProfileRequest) (*models.User, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, apperrors.NewNotFoundError("user", userID)
	}

	if req.Email != nil {
		user.Email = *req.Email
	}

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, apperrors.NewInternalError("failed to update user", err)
	}

	// Update role-specific data
	if user.Role == models.RoleMahasiswa && req.Mahasiswa != nil && user.Mahasiswa != nil {
		if req.Mahasiswa.Nama != nil {
			user.Mahasiswa.Nama = *req.Mahasiswa.Nama
		}
		if req.Mahasiswa.Prodi != nil {
			user.Mahasiswa.Prodi = *req.Mahasiswa.Prodi
		}
		if req.Mahasiswa.Angkatan != nil {
			user.Mahasiswa.Angkatan = *req.Mahasiswa.Angkatan
		}
		if err := s.repo.UpdateMahasiswa(ctx, user.Mahasiswa); err != nil {
			return nil, apperrors.NewInternalError("failed to update mahasiswa", err)
		}
	} else if user.Role == models.RoleDosen && req.Dosen != nil && user.Dosen != nil {
		if req.Dosen.Nama != nil {
			user.Dosen.Nama = *req.Dosen.Nama
		}
		if req.Dosen.Prodi != nil {
			user.Dosen.Prodi = *req.Dosen.Prodi
		}
		if err := s.repo.UpdateDosen(ctx, user.Dosen); err != nil {
			return nil, apperrors.NewInternalError("failed to update dosen", err)
		}
	} else if user.Role == models.RoleStaff && req.Staff != nil && user.Staff != nil {
		if req.Staff.Nama != nil {
			user.Staff.Nama = *req.Staff.Nama
		}
		if req.Staff.Jabatan != nil {
			user.Staff.Jabatan = *req.Staff.Jabatan
		}
		if req.Staff.Unit != nil {
			user.Staff.Unit = *req.Staff.Unit
		}
		if err := s.repo.UpdateStaff(ctx, user.Staff); err != nil {
			return nil, apperrors.NewInternalError("failed to update staff", err)
		}
	}

	// Reload user to get updated data
	return s.repo.GetUserByID(ctx, userID)
}

// GetUserByID gets user by ID
func (s *UserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

// SearchUsersRequest represents search users request
type SearchUsersRequest struct {
	Query   string `form:"q"`
	Role    string `form:"role"`
	Page    int    `form:"page,default=1"`
	PerPage int    `form:"per_page,default=20"`
}

// SearchUsers searches users
func (s *UserService) SearchUsers(ctx context.Context, req SearchUsersRequest) ([]models.User, int64, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}

	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}

	offset := (page - 1) * perPage

	var rolePtr *string
	if req.Role != "" {
		rolePtr = &req.Role
	}

	return s.repo.SearchUsers(ctx, req.Query, rolePtr, perPage, offset)
}

// GetMahasiswaByNIM gets mahasiswa by NIM
func (s *UserService) GetMahasiswaByNIM(ctx context.Context, nim string) (*models.Mahasiswa, error) {
	return s.repo.GetMahasiswaByNIM(ctx, nim)
}

// GetDosenByNIP gets dosen by NIP
func (s *UserService) GetDosenByNIP(ctx context.Context, nip string) (*models.Dosen, error) {
	return s.repo.GetDosenByNIP(ctx, nip)
}

// GetStaffByNIP gets staff by NIP
func (s *UserService) GetStaffByNIP(ctx context.Context, nip string) (*models.Staff, error) {
	return s.repo.GetStaffByNIP(ctx, nip)
}

// UploadAvatarRequest represents avatar upload request
type UploadAvatarRequest struct {
	Filename string
	Data     []byte
	MimeType string
}

// UploadAvatar uploads user avatar
func (s *UserService) UploadAvatar(ctx context.Context, userID string, req UploadAvatarRequest) (string, error) {
	// This would integrate with file storage service
	// For now, return placeholder
	return "avatar-url-placeholder", nil
}

// ListUsersRequest represents list users request
type ListUsersRequest struct {
	Role     string `form:"role"`      // Filter by role: mahasiswa, dosen, staff
	IsActive *bool  `form:"is_active"` // Filter by active status
	Page     int    `form:"page,default=1"`
	PerPage  int    `form:"per_page,default=20"`
}

// ListUsers lists all users with pagination and filters
func (s *UserService) ListUsers(ctx context.Context, req ListUsersRequest) ([]models.User, int64, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}

	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}

	offset := (page - 1) * perPage

	var rolePtr *string
	if req.Role != "" {
		rolePtr = &req.Role
	}

	return s.repo.ListUsers(ctx, rolePtr, req.IsActive, perPage, offset)
}

// AdminUpdateUserRequest represents admin update user request
type AdminUpdateUserRequest struct {
	Email     *string                 `json:"email,omitempty"`
	IsActive  *bool                   `json:"is_active,omitempty"`
	Mahasiswa *UpdateMahasiswaRequest `json:"mahasiswa,omitempty"`
	Dosen     *UpdateDosenRequest     `json:"dosen,omitempty"`
	Staff     *UpdateStaffRequest     `json:"staff,omitempty"`
}

// AdminUpdateUser updates a user (admin only)
func (s *UserService) AdminUpdateUser(ctx context.Context, userID string, req AdminUpdateUserRequest) (*models.User, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, apperrors.NewNotFoundError("user", userID)
	}

	if req.Email != nil {
		user.Email = *req.Email
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, apperrors.NewInternalError("failed to update user", err)
	}

	// Update role-specific data
	if user.Role == models.RoleMahasiswa && req.Mahasiswa != nil && user.Mahasiswa != nil {
		if req.Mahasiswa.Nama != nil {
			user.Mahasiswa.Nama = *req.Mahasiswa.Nama
		}
		if req.Mahasiswa.Prodi != nil {
			user.Mahasiswa.Prodi = *req.Mahasiswa.Prodi
		}
		if req.Mahasiswa.Angkatan != nil {
			user.Mahasiswa.Angkatan = *req.Mahasiswa.Angkatan
		}
		if err := s.repo.UpdateMahasiswa(ctx, user.Mahasiswa); err != nil {
			return nil, apperrors.NewInternalError("failed to update mahasiswa", err)
		}
	} else if user.Role == models.RoleDosen && req.Dosen != nil && user.Dosen != nil {
		if req.Dosen.Nama != nil {
			user.Dosen.Nama = *req.Dosen.Nama
		}
		if req.Dosen.Prodi != nil {
			user.Dosen.Prodi = *req.Dosen.Prodi
		}
		if err := s.repo.UpdateDosen(ctx, user.Dosen); err != nil {
			return nil, apperrors.NewInternalError("failed to update dosen", err)
		}
	} else if user.Role == models.RoleStaff && req.Staff != nil && user.Staff != nil {
		if req.Staff.Nama != nil {
			user.Staff.Nama = *req.Staff.Nama
		}
		if req.Staff.Jabatan != nil {
			user.Staff.Jabatan = *req.Staff.Jabatan
		}
		if req.Staff.Unit != nil {
			user.Staff.Unit = *req.Staff.Unit
		}
		if err := s.repo.UpdateStaff(ctx, user.Staff); err != nil {
			return nil, apperrors.NewInternalError("failed to update staff", err)
		}
	}

	// Reload user to get updated data
	return s.repo.GetUserByID(ctx, userID)
}

// DeleteUser soft deletes a user
func (s *UserService) DeleteUser(ctx context.Context, userID string) error {
	// Check if user exists
	_, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return apperrors.NewNotFoundError("user", userID)
	}

	if err := s.repo.DeleteUser(ctx, userID); err != nil {
		return apperrors.NewInternalError("failed to delete user", err)
	}

	return nil
}

// ActivateUser activates a user
func (s *UserService) ActivateUser(ctx context.Context, userID string) (*models.User, error) {
	// Check if user exists
	_, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, apperrors.NewNotFoundError("user", userID)
	}

	if err := s.repo.ActivateUser(ctx, userID); err != nil {
		return nil, apperrors.NewInternalError("failed to activate user", err)
	}

	// Reload user to get updated data
	return s.repo.GetUserByID(ctx, userID)
}

// DeactivateUser deactivates a user
func (s *UserService) DeactivateUser(ctx context.Context, userID string) (*models.User, error) {
	// Check if user exists
	_, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, apperrors.NewNotFoundError("user", userID)
	}

	if err := s.repo.DeactivateUser(ctx, userID); err != nil {
		return nil, apperrors.NewInternalError("failed to deactivate user", err)
	}

	// Reload user to get updated data
	return s.repo.GetUserByID(ctx, userID)
}

// CreateUserRequest represents admin create user request (similar to RegisterRequest)
type CreateUserRequest struct {
	Email    string          `json:"email" binding:"required,email"`
	Password string          `json:"password" binding:"required,min=8"`
	Role     models.UserRole `json:"role" binding:"required,oneof=mahasiswa dosen staff"`
	NIM      string          `json:"nim,omitempty"` // For mahasiswa
	NIP      string          `json:"nip,omitempty"` // For dosen/staff
	Nama     string          `json:"nama" binding:"required"`
	Prodi    string          `json:"prodi,omitempty"`
	Angkatan int             `json:"angkatan,omitempty"`  // For mahasiswa
	Jabatan  string          `json:"jabatan,omitempty"`   // For staff
	Unit     string          `json:"unit,omitempty"`      // For staff
	IsActive *bool           `json:"is_active,omitempty"` // Default true if not provided
}

// isConstraintViolation checks if error is a database constraint violation
func isConstraintViolation(err error) bool {
	if err == nil {
		return false
	}

	// Check for GORM duplicate entry error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	// Check for PostgreSQL unique constraint violation in error message
	errStr := strings.ToLower(err.Error())
	if strings.Contains(errStr, "unique constraint") ||
		strings.Contains(errStr, "duplicate key") ||
		strings.Contains(errStr, "violates unique constraint") ||
		strings.Contains(errStr, "23505") { // PostgreSQL unique violation error code
		return true
	}

	return false
}

// CreateUser creates a new user (admin only)
func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (*models.User, error) {
	// Check if email already exists
	existingUser, _ := s.repo.FindByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, apperrors.NewConflictError("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to hash password", err)
	}

	// Set default IsActive if not provided
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	// Use transaction to ensure atomicity
	db := s.repo.GetDB()
	var createdUserID string

	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create user within transaction
		user := &models.User{
			Email:        req.Email,
			PasswordHash: string(hashedPassword),
			Role:         req.Role,
			IsActive:     isActive,
		}

		if err := tx.Create(user).Error; err != nil {
			if isConstraintViolation(err) {
				return apperrors.NewConflictError("email already registered")
			}
			return apperrors.NewInternalError("failed to create user", err)
		}

		createdUserID = user.ID

		// Create role-specific record within transaction
		if req.Role == models.RoleMahasiswa {
			if req.NIM == "" {
				return apperrors.NewValidationError("NIM is required for mahasiswa")
			}
			// Check if NIM already exists
			var existingMahasiswa models.Mahasiswa
			if err := tx.Where("nim = ?", req.NIM).First(&existingMahasiswa).Error; err == nil {
				return apperrors.NewConflictError("NIM already registered")
			} else if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			mahasiswa := &models.Mahasiswa{
				UserID:   user.ID,
				NIM:      req.NIM,
				Nama:     req.Nama,
				Prodi:    req.Prodi,
				Angkatan: req.Angkatan,
			}
			if err := tx.Create(mahasiswa).Error; err != nil {
				if isConstraintViolation(err) {
					return apperrors.NewConflictError("NIM already registered")
				}
				return apperrors.NewInternalError("failed to create mahasiswa", err)
			}
		} else if req.Role == models.RoleDosen {
			if req.NIP == "" {
				return apperrors.NewValidationError("NIP is required for dosen")
			}
			// Check if NIP already exists
			var existingDosen models.Dosen
			if err := tx.Where("nip = ?", req.NIP).First(&existingDosen).Error; err == nil {
				return apperrors.NewConflictError("NIP already registered")
			} else if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			dosen := &models.Dosen{
				UserID: user.ID,
				NIP:    req.NIP,
				Nama:   req.Nama,
				Prodi:  req.Prodi,
			}
			if err := tx.Create(dosen).Error; err != nil {
				if isConstraintViolation(err) {
					return apperrors.NewConflictError("NIP already registered")
				}
				return apperrors.NewInternalError("failed to create dosen", err)
			}
		} else if req.Role == models.RoleStaff {
			if req.NIP == "" {
				return apperrors.NewValidationError("NIP is required for staff")
			}
			// Check if NIP already exists
			var existingStaff models.Staff
			if err := tx.Where("nip = ?", req.NIP).First(&existingStaff).Error; err == nil {
				return apperrors.NewConflictError("NIP already registered")
			} else if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			staff := &models.Staff{
				UserID:  user.ID,
				NIP:     req.NIP,
				Nama:    req.Nama,
				Jabatan: req.Jabatan,
				Unit:    req.Unit,
			}
			if err := tx.Create(staff).Error; err != nil {
				if isConstraintViolation(err) {
					return apperrors.NewConflictError("NIP already registered")
				}
				return apperrors.NewInternalError("failed to create staff", err)
			}
		}

		return nil
	})

	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			return nil, appErr
		}
		return nil, err
	}

	// Reload user to get full data with relations
	return s.repo.GetUserByID(ctx, createdUserID)
}

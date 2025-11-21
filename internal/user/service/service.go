package service

import (
	"context"

	apperrors "unsri-backend/internal/shared/errors"
	"unsri-backend/internal/shared/models"
	"unsri-backend/internal/user/repository"
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
	Email    *string `json:"email,omitempty"`
	Mahasiswa *UpdateMahasiswaRequest `json:"mahasiswa,omitempty"`
	Dosen    *UpdateDosenRequest      `json:"dosen,omitempty"`
	Staff    *UpdateStaffRequest      `json:"staff,omitempty"`
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
	Query string `form:"q"`
	Role  string `form:"role"`
	Page  int    `form:"page,default=1"`
	PerPage int  `form:"per_page,default=20"`
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


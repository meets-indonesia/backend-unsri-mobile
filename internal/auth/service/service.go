package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"
	"unsri-backend/internal/auth/repository"
	apperrors "unsri-backend/internal/shared/errors"
	"unsri-backend/internal/shared/models"
	"unsri-backend/pkg/jwt"
)

// AuthService handles authentication business logic
type AuthService struct {
	repo   *repository.AuthRepository
	jwt    *jwt.JWT
}

// NewAuthService creates a new auth service
func NewAuthService(repo *repository.AuthRepository, jwtToken *jwt.JWT) *AuthService {
	return &AuthService{
		repo: repo,
		jwt:  jwtToken,
	}
}

// LoginRequest represents login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents login response
type LoginResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	User         *UserInfo   `json:"user"`
}

// UserInfo represents user information in response
type UserInfo struct {
	ID        string                `json:"id"`
	Email     string                `json:"email"`
	Role      models.UserRole       `json:"role"`
	IsActive  bool                  `json:"is_active"`
	Mahasiswa *models.Mahasiswa     `json:"mahasiswa,omitempty"`
	Dosen     *models.Dosen         `json:"dosen,omitempty"`
	Staff     *models.Staff         `json:"staff,omitempty"`
}

// Login authenticates a user
func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperrors.NewUnauthorizedError("invalid email or password")
	}

	if !user.IsActive {
		return nil, apperrors.NewForbiddenError("account is inactive")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, apperrors.NewUnauthorizedError("invalid email or password")
	}

	accessToken, err := s.jwt.GenerateAccessToken(user.ID, string(user.Role), user.Email)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to generate access token", err)
	}

	refreshToken, err := s.jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to generate refresh token", err)
	}

	userInfo := &UserInfo{
		ID:       user.ID,
		Email:    user.Email,
		Role:     user.Role,
		IsActive: user.IsActive,
	}

	// Load role-specific data
	if user.Role == models.RoleMahasiswa && user.Mahasiswa != nil {
		userInfo.Mahasiswa = user.Mahasiswa
	} else if user.Role == models.RoleDosen && user.Dosen != nil {
		userInfo.Dosen = user.Dosen
	} else if user.Role == models.RoleStaff && user.Staff != nil {
		userInfo.Staff = user.Staff
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         userInfo,
	}, nil
}

// RegisterRequest represents registration request
type RegisterRequest struct {
	Email    string          `json:"email" binding:"required,email"`
	Password string          `json:"password" binding:"required,min=8"`
	Role     models.UserRole `json:"role" binding:"required,oneof=mahasiswa dosen staff"`
	NIM      string          `json:"nim,omitempty"` // For mahasiswa
	NIP      string          `json:"nip,omitempty"` // For dosen/staff
	Nama     string          `json:"nama" binding:"required"`
	Prodi    string          `json:"prodi,omitempty"`
	Angkatan int             `json:"angkatan,omitempty"` // For mahasiswa
	Jabatan  string          `json:"jabatan,omitempty"`  // For staff
	Unit     string          `json:"unit,omitempty"`     // For staff
}

// Register registers a new user
func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*UserInfo, error) {
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

	// Create user
	user := &models.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
		IsActive:     true,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, apperrors.NewInternalError("failed to create user", err)
	}

	// Create role-specific record
	if req.Role == models.RoleMahasiswa {
		if req.NIM == "" {
			return nil, apperrors.NewValidationError("NIM is required for mahasiswa")
		}
		mahasiswa := &models.Mahasiswa{
			UserID:   user.ID,
			NIM:      req.NIM,
			Nama:     req.Nama,
			Prodi:    req.Prodi,
			Angkatan: req.Angkatan,
		}
		if err := s.repo.CreateMahasiswa(ctx, mahasiswa); err != nil {
			return nil, apperrors.NewInternalError("failed to create mahasiswa", err)
		}
		user.Mahasiswa = mahasiswa
	} else if req.Role == models.RoleDosen {
		if req.NIP == "" {
			return nil, apperrors.NewValidationError("NIP is required for dosen")
		}
		dosen := &models.Dosen{
			UserID: user.ID,
			NIP:    req.NIP,
			Nama:   req.Nama,
			Prodi:  req.Prodi,
		}
		if err := s.repo.CreateDosen(ctx, dosen); err != nil {
			return nil, apperrors.NewInternalError("failed to create dosen", err)
		}
		user.Dosen = dosen
	} else if req.Role == models.RoleStaff {
		if req.NIP == "" {
			return nil, apperrors.NewValidationError("NIP is required for staff")
		}
		staff := &models.Staff{
			UserID: user.ID,
			NIP:    req.NIP,
			Nama:   req.Nama,
			Jabatan: req.Jabatan,
			Unit:    req.Unit,
		}
		if err := s.repo.CreateStaff(ctx, staff); err != nil {
			return nil, apperrors.NewInternalError("failed to create staff", err)
		}
		user.Staff = staff
	}

	userInfo := &UserInfo{
		ID:       user.ID,
		Email:    user.Email,
		Role:     user.Role,
		IsActive: user.IsActive,
	}

	if user.Mahasiswa != nil {
		userInfo.Mahasiswa = user.Mahasiswa
	} else if user.Dosen != nil {
		userInfo.Dosen = user.Dosen
	} else if user.Staff != nil {
		userInfo.Staff = user.Staff
	}

	return userInfo, nil
}

// RefreshTokenRequest represents refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenResponse represents refresh token response
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// RefreshToken refreshes an access token
func (s *AuthService) RefreshToken(ctx context.Context, req RefreshTokenRequest) (*RefreshTokenResponse, error) {
	claims, err := s.jwt.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, apperrors.NewUnauthorizedError("invalid refresh token")
	}

	user, err := s.repo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, apperrors.NewUnauthorizedError("user not found")
	}

	if !user.IsActive {
		return nil, apperrors.NewForbiddenError("account is inactive")
	}

	accessToken, err := s.jwt.GenerateAccessToken(user.ID, string(user.Role), user.Email)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to generate access token", err)
	}

	refreshToken, err := s.jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to generate refresh token", err)
	}

	return &RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// VerifyToken verifies a JWT token
func (s *AuthService) VerifyToken(ctx context.Context, tokenString string) (*UserInfo, error) {
	claims, err := s.jwt.ValidateToken(tokenString)
	if err != nil {
		return nil, apperrors.NewUnauthorizedError("invalid token")
	}

	user, err := s.repo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, apperrors.NewUnauthorizedError("user not found")
	}

	if !user.IsActive {
		return nil, apperrors.NewForbiddenError("account is inactive")
	}

	userInfo := &UserInfo{
		ID:       user.ID,
		Email:    user.Email,
		Role:     user.Role,
		IsActive: user.IsActive,
	}

	if user.Mahasiswa != nil {
		userInfo.Mahasiswa = user.Mahasiswa
	} else if user.Dosen != nil {
		userInfo.Dosen = user.Dosen
	} else if user.Staff != nil {
		userInfo.Staff = user.Staff
	}

	return userInfo, nil
}


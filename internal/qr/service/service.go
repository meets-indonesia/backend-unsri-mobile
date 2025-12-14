package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"time"

	"unsri-backend/internal/qr/repository"
	apperrors "unsri-backend/internal/shared/errors"
	"unsri-backend/internal/shared/models"
	userRepo "unsri-backend/internal/user/repository"
	"unsri-backend/pkg/qrcode"
)

// QRService handles QR code business logic
type QRService struct {
	repo     *repository.QRRepository
	userRepo *userRepo.UserRepository
}

// NewQRService creates a new QR service
func NewQRService(repo *repository.QRRepository, userRepo *userRepo.UserRepository) *QRService {
	return &QRService{
		repo:     repo,
		userRepo: userRepo,
	}
}

// GenerateQRRequest represents generate QR request
type GenerateQRRequest struct {
	Data     map[string]interface{} `json:"data" binding:"required"`
	Type     string                 `json:"type,omitempty"`
	Duration int                    `json:"duration,omitempty"` // minutes
}

// GenerateQRResponse represents QR generation response
type GenerateQRResponse struct {
	ID        string `json:"id"`
	QRCode    string `json:"qr_code"`
	ExpiresAt string `json:"expires_at"`
}

// GenerateQR generates a QR code
func (s *QRService) GenerateQR(ctx context.Context, createdBy string, req GenerateQRRequest) (*GenerateQRResponse, error) {
	duration := 15
	if req.Duration > 0 {
		duration = req.Duration
	}

	expiresAt := time.Now().Add(time.Duration(duration) * time.Minute)

	qrType := models.AttendanceTypeKelas
	if req.Type != "" {
		qrType = models.AttendanceType(req.Type)
	}

	dataJSON, _ := json.Marshal(req.Data)
	qrData := qrcode.QRData{
		SessionID:  "",
		ScheduleID: "",
		ExpiresAt:  expiresAt,
		Type:       string(qrType),
	}

	_, err := qrcode.GenerateQRCode(qrData)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to generate QR code", err)
	}

	session := &models.AttendanceSession{
		CreatedBy: createdBy,
		Type:      qrType,
		QRCode:    string(dataJSON),
		ExpiresAt: expiresAt,
		IsActive:  true,
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, apperrors.NewInternalError("failed to create session", err)
	}

	qrData.SessionID = session.ID
	qrImage, _ := qrcode.GenerateQRCode(qrData)

	return &GenerateQRResponse{
		ID:        session.ID,
		QRCode:    string(qrImage),
		ExpiresAt: expiresAt.Format(time.RFC3339),
	}, nil
}

// ValidateQRRequest represents validate QR request
type ValidateQRRequest struct {
	QRData string `json:"qr_data" binding:"required"`
}

// ValidateQRResponse represents QR validation response
type ValidateQRResponse struct {
	Valid   bool                   `json:"valid"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Message string                 `json:"message"`
}

// ValidateQR validates a QR code
func (s *QRService) ValidateQR(ctx context.Context, req ValidateQRRequest) (*ValidateQRResponse, error) {
	qrData, err := qrcode.ParseQRData(req.QRData)
	if err != nil {
		return &ValidateQRResponse{
			Valid:   false,
			Message: "Invalid QR code format",
		}, nil
	}

	session, err := s.repo.GetSessionByID(ctx, qrData.SessionID)
	if err != nil {
		return &ValidateQRResponse{
			Valid:   false,
			Message: "QR code session not found",
		}, nil
	}

	if !session.IsActive || time.Now().After(session.ExpiresAt) {
		return &ValidateQRResponse{
			Valid:   false,
			Message: "QR code has expired",
		}, nil
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(session.QRCode), &data); err != nil {
		return &ValidateQRResponse{
			Valid:   false,
			Message: "Failed to parse session data",
		}, nil
	}

	return &ValidateQRResponse{
		Valid:   true,
		Data:    data,
		Message: "QR code is valid",
	}, nil
}

// GetQRByID gets QR info by ID
func (s *QRService) GetQRByID(ctx context.Context, id string) (*models.AttendanceSession, error) {
	return s.repo.GetSessionByID(ctx, id)
}

// GenerateClassQRRequest represents generate class QR request
type GenerateClassQRRequest struct {
	ScheduleID string `json:"schedule_id" binding:"required"`
	Duration   int    `json:"duration,omitempty"`
}

// GenerateClassQR generates a class attendance QR code
// This QR will regenerate after each scan and attendance record
func (s *QRService) GenerateClassQR(ctx context.Context, createdBy string, req GenerateClassQRRequest) (*GenerateQRResponse, error) {
	duration := 15
	if req.Duration > 0 {
		duration = req.Duration
	}

	expiresAt := time.Now().Add(time.Duration(duration) * time.Minute)

	// Check if there's an active session for this schedule
	// If exists, deactivate it first
	existingSession, _ := s.repo.GetActiveSessionByScheduleID(ctx, req.ScheduleID)
	if existingSession != nil {
		existingSession.IsActive = false
		if err := s.repo.UpdateSession(ctx, existingSession); err != nil {
			return nil, apperrors.NewInternalError("failed to deactivate existing session", err)
		}
	}

	// Create new session
	session := &models.AttendanceSession{
		ScheduleID: &req.ScheduleID,
		CreatedBy:  createdBy,
		Type:       models.AttendanceTypeKelas,
		ExpiresAt:  expiresAt,
		IsActive:   true,
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, apperrors.NewInternalError("failed to create session", err)
	}

	qrData := qrcode.QRData{
		SessionID:  session.ID,
		ScheduleID: req.ScheduleID,
		ExpiresAt:  expiresAt,
		Type:       "kelas",
	}

	qrImage, err := qrcode.GenerateQRCode(qrData)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to generate QR code", err)
	}

	return &GenerateQRResponse{
		ID:        session.ID,
		QRCode:    string(qrImage),
		ExpiresAt: expiresAt.Format(time.RFC3339),
	}, nil
}

// RegenerateClassQR regenerates QR code for a class after scan
// This is called after attendance is recorded
func (s *QRService) RegenerateClassQR(ctx context.Context, scheduleID string, createdBy string) (*GenerateQRResponse, error) {
	// Deactivate current session
	existingSession, _ := s.repo.GetActiveSessionByScheduleID(ctx, scheduleID)
	if existingSession != nil {
		existingSession.IsActive = false
		if err := s.repo.UpdateSession(ctx, existingSession); err != nil {
			return nil, apperrors.NewInternalError("failed to deactivate existing session", err)
		}
	}

	// Generate new QR with same schedule
	req := GenerateClassQRRequest{
		ScheduleID: scheduleID,
		Duration:   15, // Default 15 minutes
	}

	return s.GenerateClassQR(ctx, createdBy, req)
}

// GenerateAccessQRRequest represents generate access QR request (deprecated, no longer needed)
type GenerateAccessQRRequest struct {
	Duration int `json:"duration,omitempty"`
}

// GenerateAccessQR generates a unique campus access QR code for user (gate access)
// This QR is unique per user and does not change
// QR contains user data for gate validation
func (s *QRService) GenerateAccessQR(ctx context.Context, userID string) (*GenerateQRResponse, error) {
	// Get user data with role-specific information
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, apperrors.NewNotFoundError("user", userID)
	}

	if !user.IsActive {
		return nil, apperrors.NewBadRequestError("user is inactive")
	}

	// Check if user already has an access QR
	userQR, err := s.repo.GetUserAccessQR(ctx, userID)
	if err == nil && userQR != nil {
		// User already has QR, regenerate with updated user data
		qrData := s.buildGateQRData(userQR.QRToken, user)
		qrImage, err := qrcode.GenerateQRCode(qrData)
		if err != nil {
			return nil, apperrors.NewInternalError("failed to generate QR code", err)
		}

		return &GenerateQRResponse{
			ID:        userQR.ID,
			QRCode:    string(qrImage),
			ExpiresAt: time.Now().Add(365 * 24 * time.Hour).Format(time.RFC3339),
		}, nil
	}

	// Generate unique token for user
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, apperrors.NewInternalError("failed to generate token", err)
	}
	qrToken := base64.URLEncoding.EncodeToString(tokenBytes)

	// Create new user access QR
	userQR = &models.UserAccessQR{
		UserID:   userID,
		QRToken:  qrToken,
		IsActive: true,
	}

	if err := s.repo.CreateUserAccessQR(ctx, userQR); err != nil {
		return nil, apperrors.NewInternalError("failed to create user access QR", err)
	}

	// Generate QR code image with user data
	qrData := s.buildGateQRData(qrToken, user)
	qrImage, err := qrcode.GenerateQRCode(qrData)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to generate QR code", err)
	}

	return &GenerateQRResponse{
		ID:        userQR.ID,
		QRCode:    string(qrImage),
		ExpiresAt: time.Now().Add(365 * 24 * time.Hour).Format(time.RFC3339),
	}, nil
}

// buildGateQRData builds QR data structure for gate access with user information
func (s *QRService) buildGateQRData(qrToken string, user *models.User) qrcode.QRData {
	qrData := qrcode.QRData{
		SessionID:  qrToken,
		ScheduleID: "",
		ExpiresAt:  time.Now().Add(365 * 24 * time.Hour), // 1 year validity
		Type:       "gate",
		UserID:     user.ID,
		QRToken:    qrToken,
		UserRole:   string(user.Role),
	}

	// Add role-specific data
	if user.Role == models.RoleMahasiswa && user.Mahasiswa != nil {
		qrData.UserName = user.Mahasiswa.Nama
		qrData.NIM = user.Mahasiswa.NIM
		qrData.Prodi = user.Mahasiswa.Prodi
	} else if user.Role == models.RoleDosen && user.Dosen != nil {
		qrData.UserName = user.Dosen.Nama
		qrData.NIP = user.Dosen.NIP
		qrData.Prodi = user.Dosen.Prodi
	} else if user.Role == models.RoleStaff && user.Staff != nil {
		qrData.UserName = user.Staff.Nama
		qrData.NIP = user.Staff.NIP
		qrData.Prodi = user.Staff.Unit // Use Unit for staff
	}

	return qrData
}

// ValidateAccessQR validates a gate access QR code
func (s *QRService) ValidateAccessQR(ctx context.Context, qrToken string) (*ValidateQRResponse, error) {
	userQR, err := s.repo.GetUserAccessQRByTokenWithUser(ctx, qrToken)
	if err != nil {
		return &ValidateQRResponse{
			Valid:   false,
			Message: "Invalid QR code",
		}, nil
	}

	if !userQR.IsActive {
		return &ValidateQRResponse{
			Valid:   false,
			Message: "QR code is inactive",
		}, nil
	}

	// Check if user is active
	if !userQR.User.IsActive {
		return &ValidateQRResponse{
			Valid:   false,
			Message: "User account is inactive",
		}, nil
	}

	// Build user data response
	userData := map[string]interface{}{
		"user_id":   userQR.UserID,
		"role":      string(userQR.User.Role),
		"is_active": userQR.User.IsActive,
	}

	// Add role-specific data
	if userQR.User.Role == models.RoleMahasiswa && userQR.User.Mahasiswa != nil {
		userData["nama"] = userQR.User.Mahasiswa.Nama
		userData["nim"] = userQR.User.Mahasiswa.NIM
		userData["prodi"] = userQR.User.Mahasiswa.Prodi
	} else if userQR.User.Role == models.RoleDosen && userQR.User.Dosen != nil {
		userData["nama"] = userQR.User.Dosen.Nama
		userData["nip"] = userQR.User.Dosen.NIP
		userData["prodi"] = userQR.User.Dosen.Prodi
	} else if userQR.User.Role == models.RoleStaff && userQR.User.Staff != nil {
		userData["nama"] = userQR.User.Staff.Nama
		userData["nip"] = userQR.User.Staff.NIP
		userData["unit"] = userQR.User.Staff.Unit
		userData["jabatan"] = userQR.User.Staff.Jabatan
	}

	return &ValidateQRResponse{
		Valid:   true,
		Message: "QR code is valid",
		Data:    userData,
	}, nil
}

// ValidateGateQRRequest represents request to validate QR from gate UNSRI
type ValidateGateQRRequest struct {
	QRData string `json:"qr_data" binding:"required"` // QR code data (JSON string from QR scan)
}

// ValidateGateQRResponse represents response for gate QR validation
type ValidateGateQRResponse struct {
	Valid    bool                   `json:"valid"`
	Allowed  bool                   `json:"allowed"`
	UserData map[string]interface{} `json:"user_data,omitempty"`
	Message  string                 `json:"message"`
}

// ValidateGateQR validates QR code from gate UNSRI
// This is a public endpoint for gate system to validate user QR codes
func (s *QRService) ValidateGateQR(ctx context.Context, req ValidateGateQRRequest) (*ValidateGateQRResponse, error) {
	// Parse QR data
	qrData, err := qrcode.ParseQRData(req.QRData)
	if err != nil {
		return &ValidateGateQRResponse{
			Valid:   false,
			Allowed: false,
			Message: "Invalid QR code format",
		}, nil
	}

	// Check if it's a gate QR
	if qrData.Type != "gate" {
		return &ValidateGateQRResponse{
			Valid:   false,
			Allowed: false,
			Message: "QR code is not a gate access QR",
		}, nil
	}

	// Validate using QR token
	validateResp, err := s.ValidateAccessQR(ctx, qrData.QRToken)
	if err != nil {
		return &ValidateGateQRResponse{
			Valid:   false,
			Allowed: false,
			Message: "Failed to validate QR code",
		}, nil
	}

	if !validateResp.Valid {
		return &ValidateGateQRResponse{
			Valid:   false,
			Allowed: false,
			Message: validateResp.Message,
		}, nil
	}

	// Verify QR data matches database
	if qrData.UserID != "" {
		userData, ok := validateResp.Data["user_id"].(string)
		if !ok || userData != qrData.UserID {
			return &ValidateGateQRResponse{
				Valid:   false,
				Allowed: false,
				Message: "QR code data mismatch",
			}, nil
		}
	}

	// User is valid and active
	return &ValidateGateQRResponse{
		Valid:    true,
		Allowed:  true,
		UserData: validateResp.Data,
		Message:  "Access granted",
	}, nil
}

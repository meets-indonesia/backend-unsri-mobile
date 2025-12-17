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

	"github.com/google/uuid"
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

	// Generate new session ID (UUID) - always create new session
	sessionID := uuid.New().String()

	// Generate unique token for user (legacy, kept for backward compatibility)
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, apperrors.NewInternalError("failed to generate token", err)
	}
	qrToken := base64.URLEncoding.EncodeToString(tokenBytes)

	// Create new user access QR with new session ID
	userQR := &models.UserAccessQR{
		UserID:    userID,
		SessionID: sessionID,
		QRToken:   qrToken,
		IsActive:  true, // Start as active (tap-in ready)
		ExpiresAt: nil,  // Not expired yet
	}

	if err := s.repo.CreateUserAccessQR(ctx, userQR); err != nil {
		return nil, apperrors.NewInternalError("failed to create user access QR", err)
	}

	// Generate QR code image with user data (only required fields)
	qrData := s.buildGateQRData(sessionID, user)
	qrImage, err := qrcode.GenerateQRCode(qrData)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to generate QR code", err)
	}

	// Encode PNG bytes to base64 string
	qrCodeBase64 := base64.StdEncoding.EncodeToString(qrImage)

	return &GenerateQRResponse{
		ID:        userQR.ID,
		QRCode:    qrCodeBase64,
		ExpiresAt: "", // No expiration until tap-out
	}, nil
}

// buildGateQRData builds QR data structure for gate access with user information
// Only includes: session_id, user_role, user_name, nim/nip, prodi
func (s *QRService) buildGateQRData(sessionID string, user *models.User) qrcode.QRData {
	qrData := qrcode.QRData{
		SessionID: sessionID,
		Type:      "gate",
		UserRole:  string(user.Role),
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

// ValidateAccessQR validates a gate access QR code using session_id
// Implements tap-in/tap-out logic:
// - First scan (tap-in/masuk): is_active = true, allow access
// - Second scan (tap-out/keluar): is_active = false, expires_at = now, allow access
// - After tap-out: session expired, must generate new QR
func (s *QRService) ValidateAccessQR(ctx context.Context, sessionID string) (*ValidateQRResponse, error) {
	// Get user access QR by session ID (only non-expired sessions)
	userQR, err := s.repo.GetUserAccessQRBySessionID(ctx, sessionID)
	if err != nil {
		return &ValidateQRResponse{
			Valid:   false,
			Message: "Invalid QR code or session expired",
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
		"user_id":    userQR.UserID,
		"session_id": userQR.SessionID,
		"role":       string(userQR.User.Role),
		"is_active":  userQR.User.IsActive,
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

	now := time.Now()

	// Determine if this is tap-in (first scan) or tap-out (second scan)
	// Use a more robust approach: check if UpdatedAt is within a small window of CreatedAt
	// This indicates the QR hasn't been tapped in yet (tap-in state)
	// Use a 5-second window to account for database timestamp precision and small delays
	timeWindow := 5 * time.Second
	timeDiff := userQR.UpdatedAt.Sub(userQR.CreatedAt)
	isFirstScan := timeDiff >= 0 && timeDiff <= timeWindow

	if isFirstScan && userQR.IsActive && userQR.ExpiresAt == nil {
		// First scan: tap-in (masuk) - keep active, allow access
		// Use atomic update with WHERE clause to ensure only one tap-in can succeed
		// This prevents race conditions when multiple scans happen simultaneously
		userData["action"] = "tap_in"
		userData["status"] = "masuk"
		userData["message"] = "Access granted (tap-in). Status QR aktif. Next scan will be tap-out."

		// Atomic update: only update if UpdatedAt hasn't changed significantly (still in tap-in state)
		// This ensures only the first scan can mark as tap-in
		updateTime := now
		rowsAffected, err := s.repo.UpdateUserAccessQRAtomic(ctx, userQR.SessionID, updateTime, timeWindow)
		if err != nil {
			return &ValidateQRResponse{
				Valid:   false,
				Message: "Failed to update QR status",
			}, nil
		}

		// If no rows were affected, it means another request already processed tap-in
		// This can happen in high concurrency scenarios
		if rowsAffected == 0 {
			// Re-fetch to get updated state
			updatedQR, err := s.repo.GetUserAccessQRBySessionID(ctx, userQR.SessionID)
			if err != nil {
				return &ValidateQRResponse{
					Valid:   false,
					Message: "Failed to verify QR status",
				}, nil
			}

			// Check if it's now in tap-out state (UpdatedAt was updated by another request)
			timeDiffAfter := updatedQR.UpdatedAt.Sub(updatedQR.CreatedAt)
			if timeDiffAfter > timeWindow && updatedQR.IsActive && updatedQR.ExpiresAt == nil {
				// Another request already processed tap-in, this should be tap-out
				userData["action"] = "tap_out"
				userData["status"] = "keluar"
				userData["message"] = "Access granted (tap-out). Session expired. Generate new QR for next entry."

				// Set expired and inactive
				updatedQR.IsActive = false
				updatedQR.ExpiresAt = &now

				if err := s.repo.UpdateUserAccessQR(ctx, updatedQR); err != nil {
					return &ValidateQRResponse{
						Valid:   false,
						Message: "Failed to update QR status",
					}, nil
				}
			} else {
				// Still in tap-in state or invalid state
				return &ValidateQRResponse{
					Valid:   false,
					Message: "QR status conflict, please try again",
				}, nil
			}
		} else {
			// Successfully updated, update local object
			userQR.UpdatedAt = updateTime
		}
	} else if !isFirstScan && userQR.IsActive && userQR.ExpiresAt == nil {
		// Second scan: tap-out (keluar) - set inactive, expire session, allow access
		userData["action"] = "tap_out"
		userData["status"] = "keluar"
		userData["message"] = "Access granted (tap-out). Session expired. Generate new QR for next entry."

		// Set expired and inactive
		userQR.IsActive = false
		userQR.ExpiresAt = &now

		// Update in database
		if err := s.repo.UpdateUserAccessQR(ctx, userQR); err != nil {
			return &ValidateQRResponse{
				Valid:   false,
				Message: "Failed to update QR status",
			}, nil
		}
	} else {
		// Invalid state
		return &ValidateQRResponse{
			Valid:   false,
			Message: "Invalid QR state or already used",
		}, nil
	}

	return &ValidateQRResponse{
		Valid:   true,
		Data:    userData,
		Message: userData["message"].(string),
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

	// Validate using session_id from QR data
	if qrData.SessionID == "" {
		return &ValidateGateQRResponse{
			Valid:   false,
			Allowed: false,
			Message: "Session ID is required in QR code",
		}, nil
	}

	validateResp, err := s.ValidateAccessQR(ctx, qrData.SessionID)
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

	// Verify QR data matches database (optional validation)
	if qrData.UserRole != "" {
		userRole, ok := validateResp.Data["role"].(string)
		if !ok || userRole != qrData.UserRole {
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

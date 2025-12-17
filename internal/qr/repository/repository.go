package repository

import (
	"context"
	"errors"
	"time"

	"unsri-backend/internal/shared/models"

	"gorm.io/gorm"
)

// QRRepository handles QR code data operations
type QRRepository struct {
	db *gorm.DB
}

// NewQRRepository creates a new QR repository
func NewQRRepository(db *gorm.DB) *QRRepository {
	return &QRRepository{db: db}
}

// CreateSession creates a new QR session
func (r *QRRepository) CreateSession(ctx context.Context, session *models.AttendanceSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

// GetSessionByID gets a session by ID
func (r *QRRepository) GetSessionByID(ctx context.Context, id string) (*models.AttendanceSession, error) {
	var session models.AttendanceSession
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("session not found")
		}
		return nil, err
	}
	return &session, nil
}

// UpdateSession updates a session
func (r *QRRepository) UpdateSession(ctx context.Context, session *models.AttendanceSession) error {
	return r.db.WithContext(ctx).Save(session).Error
}

// GetActiveSessions gets active sessions
func (r *QRRepository) GetActiveSessions(ctx context.Context, sessionType *string) ([]models.AttendanceSession, error) {
	var sessions []models.AttendanceSession
	query := r.db.WithContext(ctx).Where("is_active = ? AND expires_at > ?", true, time.Now())

	if sessionType != nil {
		query = query.Where("type = ?", *sessionType)
	}

	if err := query.Order("created_at DESC").Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

// GetActiveSessionByScheduleID gets active session for a schedule
func (r *QRRepository) GetActiveSessionByScheduleID(ctx context.Context, scheduleID string) (*models.AttendanceSession, error) {
	var session models.AttendanceSession
	if err := r.db.WithContext(ctx).
		Where("schedule_id = ? AND is_active = ? AND expires_at > ?", scheduleID, true, time.Now()).
		Order("created_at DESC").
		First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no active session found")
		}
		return nil, err
	}
	return &session, nil
}

// GetUserAccessQR gets user access QR by user ID
func (r *QRRepository) GetUserAccessQR(ctx context.Context, userID string) (*models.UserAccessQR, error) {
	var userQR models.UserAccessQR
	if err := r.db.WithContext(ctx).Where("user_id = ? AND is_active = ?", userID, true).First(&userQR).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user access QR not found")
		}
		return nil, err
	}
	return &userQR, nil
}

// CreateUserAccessQR creates a user access QR
func (r *QRRepository) CreateUserAccessQR(ctx context.Context, userQR *models.UserAccessQR) error {
	return r.db.WithContext(ctx).Create(userQR).Error
}

// UpdateUserAccessQR updates a user access QR
func (r *QRRepository) UpdateUserAccessQR(ctx context.Context, userQR *models.UserAccessQR) error {
	return r.db.WithContext(ctx).Save(userQR).Error
}

// UpdateUserAccessQRAtomic atomically updates UpdatedAt for tap-in detection
// Only updates if UpdatedAt is still within the time window from CreatedAt (tap-in state)
// Returns the number of rows affected (0 if condition not met, 1 if updated)
func (r *QRRepository) UpdateUserAccessQRAtomic(ctx context.Context, sessionID string, updateTime time.Time, timeWindow time.Duration) (int64, error) {
	result := r.db.WithContext(ctx).
		Model(&models.UserAccessQR{}).
		Where("session_id = ? AND is_active = ? AND expires_at IS NULL", sessionID, true).
		Where("EXTRACT(EPOCH FROM (updated_at - created_at)) <= ?", timeWindow.Seconds()).
		Update("updated_at", updateTime)

	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// GetUserAccessQRByToken gets user access QR by token
func (r *QRRepository) GetUserAccessQRByToken(ctx context.Context, token string) (*models.UserAccessQR, error) {
	var userQR models.UserAccessQR
	if err := r.db.WithContext(ctx).Where("qr_token = ? AND is_active = ?", token, true).First(&userQR).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user access QR not found")
		}
		return nil, err
	}
	return &userQR, nil
}

// GetUserAccessQRByTokenWithUser gets user access QR by token with user data
func (r *QRRepository) GetUserAccessQRByTokenWithUser(ctx context.Context, token string) (*models.UserAccessQR, error) {
	var userQR models.UserAccessQR
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("User.Mahasiswa").
		Preload("User.Dosen").
		Preload("User.Staff").
		Where("qr_token = ? AND is_active = ?", token, true).
		First(&userQR).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user access QR not found")
		}
		return nil, err
	}
	return &userQR, nil
}

// GetUserAccessQRBySessionID gets user access QR by session ID with user data
// Returns QR that is not expired (for tap-in) or active QR (for tap-out check)
func (r *QRRepository) GetUserAccessQRBySessionID(ctx context.Context, sessionID string) (*models.UserAccessQR, error) {
	var userQR models.UserAccessQR
	query := r.db.WithContext(ctx).
		Preload("User").
		Preload("User.Mahasiswa").
		Preload("User.Dosen").
		Preload("User.Staff").
		Where("session_id = ?", sessionID)

	// Get active QR (not expired yet) - for tap-in
	// Or get QR that is_active = true (for tap-out, even if expires_at is set)
	query = query.Where("(expires_at IS NULL OR expires_at > ?) AND is_active = ?", time.Now(), true)

	if err := query.First(&userQR).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user access QR not found or expired")
		}
		return nil, err
	}
	return &userQR, nil
}

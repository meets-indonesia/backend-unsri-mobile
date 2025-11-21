package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"unsri-backend/internal/shared/models"
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


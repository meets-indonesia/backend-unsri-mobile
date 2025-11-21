package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"unsri-backend/internal/shared/models"
)

// QuickActionsRepository handles quick actions data operations
type QuickActionsRepository struct {
	db *gorm.DB
}

// NewQuickActionsRepository creates a new quick actions repository
func NewQuickActionsRepository(db *gorm.DB) *QuickActionsRepository {
	return &QuickActionsRepository{db: db}
}

// GetTranscriptByStudentID gets transcript for a student
func (r *QuickActionsRepository) GetTranscriptByStudentID(ctx context.Context, studentID string) (*models.Transcript, error) {
	var transcript models.Transcript
	if err := r.db.WithContext(ctx).
		Where("student_id = ?", studentID).
		Order("semester DESC").
		First(&transcript).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transcript not found")
		}
		return nil, err
	}
	return &transcript, nil
}

// GetKRSByStudentID gets KRS for a student
func (r *QuickActionsRepository) GetKRSByStudentID(ctx context.Context, studentID string, semester *string) (*models.KRS, error) {
	var krs models.KRS
	query := r.db.WithContext(ctx).Where("student_id = ?", studentID)

	if semester != nil {
		query = query.Where("semester = ?", *semester)
	}

	if err := query.Order("semester DESC").First(&krs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("KRS not found")
		}
		return nil, err
	}
	return &krs, nil
}

// GetBimbingansByStudentID gets bimbingan records for a student
func (r *QuickActionsRepository) GetBimbingansByStudentID(ctx context.Context, studentID string, limit int) ([]models.Bimbingan, error) {
	var bimbingans []models.Bimbingan
	query := r.db.WithContext(ctx).Where("student_id = ?", studentID)

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Order("date DESC").Find(&bimbingans).Error; err != nil {
		return nil, err
	}
	return bimbingans, nil
}

// GetBimbingansByDosenID gets bimbingan records for a dosen
func (r *QuickActionsRepository) GetBimbingansByDosenID(ctx context.Context, dosenID string, limit int) ([]models.Bimbingan, error) {
	var bimbingans []models.Bimbingan
	query := r.db.WithContext(ctx).Where("dosen_id = ?", dosenID)

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Order("date DESC").Find(&bimbingans).Error; err != nil {
		return nil, err
	}
	return bimbingans, nil
}


package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"unsri-backend/internal/shared/models"
)

// FileRepository handles file data operations
type FileRepository struct {
	db *gorm.DB
}

// NewFileRepository creates a new file repository
func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{db: db}
}

// CreateFile creates a file record
func (r *FileRepository) CreateFile(ctx context.Context, file *models.File) error {
	return r.db.WithContext(ctx).Create(file).Error
}

// GetFileByID gets a file by ID
func (r *FileRepository) GetFileByID(ctx context.Context, id string) (*models.File, error) {
	var file models.File
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&file).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("file not found")
		}
		return nil, err
	}
	return &file, nil
}

// GetFilesByUserID gets files for a user
func (r *FileRepository) GetFilesByUserID(ctx context.Context, userID string, fileType *string, limit, offset int) ([]models.File, int64, error) {
	var files []models.File
	var total int64

	query := r.db.WithContext(ctx).Model(&models.File{}).Where("user_id = ?", userID)

	if fileType != nil {
		query = query.Where("file_type = ?", *fileType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&files).Error; err != nil {
		return nil, 0, err
	}

	return files, total, nil
}

// UpdateFile updates a file record
func (r *FileRepository) UpdateFile(ctx context.Context, file *models.File) error {
	return r.db.WithContext(ctx).Save(file).Error
}

// DeleteFile soft deletes a file
func (r *FileRepository) DeleteFile(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.File{}, "id = ?", id).Error
}

// GetAvatarByUserID gets user avatar
func (r *FileRepository) GetAvatarByUserID(ctx context.Context, userID string) (*models.File, error) {
	var file models.File
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND file_type = ?", userID, "avatar").
		Order("created_at DESC").
		First(&file).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &file, nil
}


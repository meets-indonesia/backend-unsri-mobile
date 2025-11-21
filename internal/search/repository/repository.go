package repository

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"unsri-backend/internal/shared/models"
)

// SearchRepository handles search data operations
type SearchRepository struct {
	db *gorm.DB
}

// NewSearchRepository creates a new search repository
func NewSearchRepository(db *gorm.DB) *SearchRepository {
	return &SearchRepository{db: db}
}

// SearchUsers searches for users
func (r *SearchRepository) SearchUsers(ctx context.Context, query string, role *string, limit, offset int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&models.User{})

	// Build search query
	searchTerm := "%" + strings.ToLower(query) + "%"
	dbQuery = dbQuery.Where(
		"LOWER(name) LIKE ? OR LOWER(email) LIKE ? OR LOWER(nim) LIKE ? OR LOWER(nip) LIKE ?",
		searchTerm, searchTerm, searchTerm, searchTerm,
	)

	if role != nil && *role != "" {
		dbQuery = dbQuery.Where("role = ?", *role)
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := dbQuery.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// SearchCourses searches for courses
func (r *SearchRepository) SearchCourses(ctx context.Context, query string, limit, offset int) ([]models.Course, int64, error) {
	var courses []models.Course
	var total int64

	searchTerm := "%" + strings.ToLower(query) + "%"
	dbQuery := r.db.WithContext(ctx).Model(&models.Course{}).
		Where("LOWER(name) LIKE ? OR LOWER(code) LIKE ? OR LOWER(description) LIKE ?",
			searchTerm, searchTerm, searchTerm)

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := dbQuery.Limit(limit).Offset(offset).Find(&courses).Error; err != nil {
		return nil, 0, err
	}

	return courses, total, nil
}

// SearchSchedules searches for schedules
func (r *SearchRepository) SearchSchedules(ctx context.Context, query string, userID *string, limit, offset int) ([]models.Schedule, int64, error) {
	var schedules []models.Schedule
	var total int64

	searchTerm := "%" + strings.ToLower(query) + "%"
	dbQuery := r.db.WithContext(ctx).Model(&models.Schedule{}).
		Joins("LEFT JOIN courses ON schedules.course_id = courses.id").
		Where("LOWER(courses.name) LIKE ? OR LOWER(courses.code) LIKE ? OR LOWER(schedules.room) LIKE ?",
			searchTerm, searchTerm, searchTerm)

	if userID != nil {
		dbQuery = dbQuery.Where("schedules.dosen_id = ? OR schedules.id IN (SELECT schedule_id FROM schedule_students WHERE student_id = ?)",
			*userID, *userID)
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := dbQuery.Limit(limit).Offset(offset).Find(&schedules).Error; err != nil {
		return nil, 0, err
	}

	return schedules, total, nil
}

// SearchGlobal performs global search across multiple entities
func (r *SearchRepository) SearchGlobal(ctx context.Context, query string, types []string, limit int) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	searchTerm := "%" + strings.ToLower(query) + "%"

	for _, searchType := range types {
		switch searchType {
		case "users":
			var users []models.User
			r.db.WithContext(ctx).Model(&models.User{}).
				Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ? OR LOWER(nim) LIKE ? OR LOWER(nip) LIKE ?",
					searchTerm, searchTerm, searchTerm, searchTerm).
				Limit(limit).
				Find(&users)
			result["users"] = users

		case "courses":
			var courses []models.Course
			r.db.WithContext(ctx).Model(&models.Course{}).
				Where("LOWER(name) LIKE ? OR LOWER(code) LIKE ?",
					searchTerm, searchTerm).
				Limit(limit).
				Find(&courses)
			result["courses"] = courses

		case "schedules":
			var schedules []models.Schedule
			r.db.WithContext(ctx).Model(&models.Schedule{}).
				Joins("LEFT JOIN courses ON schedules.course_id = courses.id").
				Where("LOWER(courses.name) LIKE ? OR LOWER(courses.code) LIKE ?",
					searchTerm, searchTerm).
				Limit(limit).
				Find(&schedules)
			result["schedules"] = schedules

		case "broadcasts":
			var broadcasts []models.Broadcast
			r.db.WithContext(ctx).Model(&models.Broadcast{}).
				Where("LOWER(title) LIKE ? OR LOWER(content) LIKE ?",
					searchTerm, searchTerm).
				Limit(limit).
				Find(&broadcasts)
			result["broadcasts"] = broadcasts
		}
	}

	return result, nil
}

// SearchByType performs type-specific search
func (r *SearchRepository) SearchByType(ctx context.Context, searchType, query string, filters map[string]interface{}, limit, offset int) (interface{}, int64, error) {
	searchTerm := "%" + strings.ToLower(query) + "%"

	switch searchType {
	case "users":
		var users []models.User
		var total int64
		dbQuery := r.db.WithContext(ctx).Model(&models.User{}).
			Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ? OR LOWER(nim) LIKE ? OR LOWER(nip) LIKE ?",
				searchTerm, searchTerm, searchTerm, searchTerm)

		if role, ok := filters["role"].(string); ok && role != "" {
			dbQuery = dbQuery.Where("role = ?", role)
		}

		dbQuery.Count(&total)
		dbQuery.Limit(limit).Offset(offset).Find(&users)
		return users, total, nil

	case "courses":
		var courses []models.Course
		var total int64
		dbQuery := r.db.WithContext(ctx).Model(&models.Course{}).
			Where("LOWER(name) LIKE ? OR LOWER(code) LIKE ? OR LOWER(description) LIKE ?",
				searchTerm, searchTerm, searchTerm)

		if semester, ok := filters["semester"].(string); ok && semester != "" {
			dbQuery = dbQuery.Where("semester = ?", semester)
		}

		dbQuery.Count(&total)
		dbQuery.Limit(limit).Offset(offset).Find(&courses)
		return courses, total, nil

	case "schedules":
		var schedules []models.Schedule
		var total int64
		dbQuery := r.db.WithContext(ctx).Model(&models.Schedule{}).
			Joins("LEFT JOIN courses ON schedules.course_id = courses.id").
			Where("LOWER(courses.name) LIKE ? OR LOWER(courses.code) LIKE ? OR LOWER(schedules.room) LIKE ?",
				searchTerm, searchTerm, searchTerm)

		if userID, ok := filters["user_id"].(string); ok && userID != "" {
			dbQuery = dbQuery.Where("schedules.dosen_id = ? OR schedules.id IN (SELECT schedule_id FROM schedule_students WHERE student_id = ?)",
				userID, userID)
		}

		dbQuery.Count(&total)
		dbQuery.Limit(limit).Offset(offset).Find(&schedules)
		return schedules, total, nil

	default:
		return nil, 0, fmt.Errorf("unsupported search type: %s", searchType)
	}
}


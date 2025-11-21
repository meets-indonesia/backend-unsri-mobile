package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"unsri-backend/internal/shared/models"
)

// CourseRepository handles course data operations
type CourseRepository struct {
	db *gorm.DB
}

// NewCourseRepository creates a new course repository
func NewCourseRepository(db *gorm.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

// CreateCourse creates a new course
func (r *CourseRepository) CreateCourse(ctx context.Context, course *models.Course) error {
	return r.db.WithContext(ctx).Create(course).Error
}

// GetCourseByID gets a course by ID
func (r *CourseRepository) GetCourseByID(ctx context.Context, id string) (*models.Course, error) {
	var course models.Course
	if err := r.db.WithContext(ctx).Preload("Classes").Where("id = ?", id).First(&course).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("course not found")
		}
		return nil, err
	}
	return &course, nil
}

// GetCourseByCode gets a course by code
func (r *CourseRepository) GetCourseByCode(ctx context.Context, code string) (*models.Course, error) {
	var course models.Course
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&course).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("course not found")
		}
		return nil, err
	}
	return &course, nil
}

// GetAllCourses gets all courses with filters
func (r *CourseRepository) GetAllCourses(ctx context.Context, prodi *string, isActive *bool, limit, offset int) ([]models.Course, int64, error) {
	var courses []models.Course
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Course{})

	if prodi != nil {
		query = query.Where("prodi = ?", *prodi)
	}
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Classes").Limit(limit).Offset(offset).Order("code ASC").Find(&courses).Error; err != nil {
		return nil, 0, err
	}

	return courses, total, nil
}

// UpdateCourse updates a course
func (r *CourseRepository) UpdateCourse(ctx context.Context, course *models.Course) error {
	return r.db.WithContext(ctx).Save(course).Error
}

// DeleteCourse soft deletes a course
func (r *CourseRepository) DeleteCourse(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Course{}, "id = ?", id).Error
}

// CreateClass creates a new class
func (r *CourseRepository) CreateClass(ctx context.Context, class *models.Class) error {
	return r.db.WithContext(ctx).Create(class).Error
}

// GetClassByID gets a class by ID
func (r *CourseRepository) GetClassByID(ctx context.Context, id string) (*models.Class, error) {
	var class models.Class
	if err := r.db.WithContext(ctx).Preload("Course").Preload("Enrollments").Where("id = ?", id).First(&class).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("class not found")
		}
		return nil, err
	}
	return &class, nil
}

// GetAllClasses gets all classes with filters
func (r *CourseRepository) GetAllClasses(ctx context.Context, courseID *string, dosenID *string, semester *string, limit, offset int) ([]models.Class, int64, error) {
	var classes []models.Class
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Class{})

	if courseID != nil {
		query = query.Where("course_id = ?", *courseID)
	}
	if dosenID != nil {
		query = query.Where("dosen_id = ?", *dosenID)
	}
	if semester != nil {
		query = query.Where("semester = ?", *semester)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Course").Limit(limit).Offset(offset).Order("class_code ASC").Find(&classes).Error; err != nil {
		return nil, 0, err
	}

	return classes, total, nil
}

// UpdateClass updates a class
func (r *CourseRepository) UpdateClass(ctx context.Context, class *models.Class) error {
	return r.db.WithContext(ctx).Save(class).Error
}

// DeleteClass soft deletes a class
func (r *CourseRepository) DeleteClass(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Class{}, "id = ?", id).Error
}

// GetClassesByStudentID gets classes enrolled by a student
func (r *CourseRepository) GetClassesByStudentID(ctx context.Context, studentID string) ([]models.Class, error) {
	var classes []models.Class
	if err := r.db.WithContext(ctx).
		Joins("JOIN enrollments ON enrollments.class_id = classes.id").
		Where("enrollments.student_id = ? AND enrollments.status = ?", studentID, "active").
		Preload("Course").
		Find(&classes).Error; err != nil {
		return nil, err
	}
	return classes, nil
}

// GetClassesByLecturerID gets classes taught by a lecturer
func (r *CourseRepository) GetClassesByLecturerID(ctx context.Context, lecturerID string) ([]models.Class, error) {
	var classes []models.Class
	if err := r.db.WithContext(ctx).
		Where("dosen_id = ? OR assistant_dosen_id = ?", lecturerID, lecturerID).
		Preload("Course").
		Find(&classes).Error; err != nil {
		return nil, err
	}
	return classes, nil
}

// CreateEnrollment creates a new enrollment
func (r *CourseRepository) CreateEnrollment(ctx context.Context, enrollment *models.Enrollment) error {
	return r.db.WithContext(ctx).Create(enrollment).Error
}

// GetEnrollmentByID gets an enrollment by ID
func (r *CourseRepository) GetEnrollmentByID(ctx context.Context, id string) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	if err := r.db.WithContext(ctx).Preload("Student").Preload("Class").Where("id = ?", id).First(&enrollment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("enrollment not found")
		}
		return nil, err
	}
	return &enrollment, nil
}

// GetEnrollmentsByStudentID gets enrollments for a student
func (r *CourseRepository) GetEnrollmentsByStudentID(ctx context.Context, studentID string) ([]models.Enrollment, error) {
	var enrollments []models.Enrollment
	if err := r.db.WithContext(ctx).
		Where("student_id = ?", studentID).
		Preload("Class").Preload("Class.Course").
		Order("enrollment_date DESC").
		Find(&enrollments).Error; err != nil {
		return nil, err
	}
	return enrollments, nil
}

// UpdateEnrollment updates an enrollment
func (r *CourseRepository) UpdateEnrollment(ctx context.Context, enrollment *models.Enrollment) error {
	return r.db.WithContext(ctx).Save(enrollment).Error
}


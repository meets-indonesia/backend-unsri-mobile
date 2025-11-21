package service

import (
	"context"
	"time"

	apperrors "unsri-backend/internal/shared/errors"
	"unsri-backend/internal/shared/models"
	"unsri-backend/internal/course/repository"
)

// CourseService handles course business logic
type CourseService struct {
	repo *repository.CourseRepository
}

// NewCourseService creates a new course service
func NewCourseService(repo *repository.CourseRepository) *CourseService {
	return &CourseService{repo: repo}
}

// CreateCourseRequest represents create course request
type CreateCourseRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	NameEn      string `json:"name_en,omitempty"`
	Credits     int    `json:"credits" binding:"required"`
	Semester    int    `json:"semester,omitempty"`
	Prodi       string `json:"prodi,omitempty"`
	Description string `json:"description,omitempty"`
}

// CreateCourse creates a new course
func (s *CourseService) CreateCourse(ctx context.Context, req CreateCourseRequest) (*models.Course, error) {
	course := &models.Course{
		Code:        req.Code,
		Name:        req.Name,
		NameEn:      req.NameEn,
		Credits:     req.Credits,
		Semester:    req.Semester,
		Prodi:       req.Prodi,
		Description: req.Description,
		IsActive:    true,
	}

	if err := s.repo.CreateCourse(ctx, course); err != nil {
		return nil, apperrors.NewInternalError("failed to create course", err)
	}

	return course, nil
}

// GetCourseByID gets a course by ID
func (s *CourseService) GetCourseByID(ctx context.Context, id string) (*models.Course, error) {
	return s.repo.GetCourseByID(ctx, id)
}

// GetCoursesRequest represents get courses request
type GetCoursesRequest struct {
	Prodi    string `form:"prodi"`
	IsActive *bool  `form:"is_active"`
	Page     int    `form:"page,default=1"`
	PerPage  int    `form:"per_page,default=20"`
}

// GetCourses gets all courses
func (s *CourseService) GetCourses(ctx context.Context, req GetCoursesRequest) ([]models.Course, int64, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}

	var prodiPtr *string
	if req.Prodi != "" {
		prodiPtr = &req.Prodi
	}

	return s.repo.GetAllCourses(ctx, prodiPtr, req.IsActive, perPage, (page-1)*perPage)
}

// UpdateCourseRequest represents update course request
type UpdateCourseRequest struct {
	Name        *string `json:"name,omitempty"`
	NameEn      *string `json:"name_en,omitempty"`
	Credits     *int    `json:"credits,omitempty"`
	Semester    *int    `json:"semester,omitempty"`
	Prodi       *string `json:"prodi,omitempty"`
	Description *string `json:"description,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

// UpdateCourse updates a course
func (s *CourseService) UpdateCourse(ctx context.Context, id string, req UpdateCourseRequest) (*models.Course, error) {
	course, err := s.repo.GetCourseByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewNotFoundError("course", id)
	}

	if req.Name != nil {
		course.Name = *req.Name
	}
	if req.NameEn != nil {
		course.NameEn = *req.NameEn
	}
	if req.Credits != nil {
		course.Credits = *req.Credits
	}
	if req.Semester != nil {
		course.Semester = *req.Semester
	}
	if req.Prodi != nil {
		course.Prodi = *req.Prodi
	}
	if req.Description != nil {
		course.Description = *req.Description
	}
	if req.IsActive != nil {
		course.IsActive = *req.IsActive
	}

	if err := s.repo.UpdateCourse(ctx, course); err != nil {
		return nil, apperrors.NewInternalError("failed to update course", err)
	}

	return course, nil
}

// DeleteCourse deletes a course
func (s *CourseService) DeleteCourse(ctx context.Context, id string) error {
	_, err := s.repo.GetCourseByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFoundError("course", id)
	}
	return s.repo.DeleteCourse(ctx, id)
}

// CreateClassRequest represents create class request
type CreateClassRequest struct {
	CourseID         string  `json:"course_id" binding:"required"`
	ClassCode        string  `json:"class_code" binding:"required"`
	ClassName        string  `json:"class_name,omitempty"`
	Semester         string  `json:"semester" binding:"required"`
	AcademicYear     string  `json:"academic_year,omitempty"`
	Capacity         int     `json:"capacity,omitempty"`
	DosenID          string  `json:"dosen_id" binding:"required"`
	AssistantDosenID *string `json:"assistant_dosen_id,omitempty"`
	Room             string  `json:"room,omitempty"`
	DayOfWeek        int     `json:"day_of_week" binding:"required,min=0,max=6"`
	StartTime        string  `json:"start_time" binding:"required"`
	EndTime          string  `json:"end_time" binding:"required"`
}

// CreateClass creates a new class
func (s *CourseService) CreateClass(ctx context.Context, req CreateClassRequest) (*models.Class, error) {
	// Parse times
	startTime, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		return nil, apperrors.NewValidationError("invalid start_time format, use HH:MM")
	}

	endTime, err := time.Parse("15:04", req.EndTime)
	if err != nil {
		return nil, apperrors.NewValidationError("invalid end_time format, use HH:MM")
	}

	// Use current date for time fields (will be used with day_of_week)
	now := time.Now()
	startDateTime := time.Date(now.Year(), now.Month(), now.Day(), startTime.Hour(), startTime.Minute(), 0, 0, now.Location())
	endDateTime := time.Date(now.Year(), now.Month(), now.Day(), endTime.Hour(), endTime.Minute(), 0, 0, now.Location())

	class := &models.Class{
		CourseID:         req.CourseID,
		ClassCode:        req.ClassCode,
		ClassName:        req.ClassName,
		Semester:         req.Semester,
		AcademicYear:     req.AcademicYear,
		Capacity:         req.Capacity,
		DosenID:          req.DosenID,
		AssistantDosenID: req.AssistantDosenID,
		Room:             req.Room,
		DayOfWeek:        req.DayOfWeek,
		StartTime:        startDateTime,
		EndTime:          endDateTime,
		IsActive:         true,
		Enrolled:         0,
	}

	if err := s.repo.CreateClass(ctx, class); err != nil {
		return nil, apperrors.NewInternalError("failed to create class", err)
	}

	return class, nil
}

// GetClassByID gets a class by ID
func (s *CourseService) GetClassByID(ctx context.Context, id string) (*models.Class, error) {
	return s.repo.GetClassByID(ctx, id)
}

// GetClassesRequest represents get classes request
type GetClassesRequest struct {
	CourseID string `form:"course_id"`
	DosenID  string `form:"dosen_id"`
	Semester string `form:"semester"`
	Page     int    `form:"page,default=1"`
	PerPage  int    `form:"per_page,default=20"`
}

// GetClasses gets all classes
func (s *CourseService) GetClasses(ctx context.Context, req GetClassesRequest) ([]models.Class, int64, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}

	var courseIDPtr, dosenIDPtr, semesterPtr *string
	if req.CourseID != "" {
		courseIDPtr = &req.CourseID
	}
	if req.DosenID != "" {
		dosenIDPtr = &req.DosenID
	}
	if req.Semester != "" {
		semesterPtr = &req.Semester
	}

	return s.repo.GetAllClasses(ctx, courseIDPtr, dosenIDPtr, semesterPtr, perPage, (page-1)*perPage)
}

// GetClassesByStudent gets classes for a student
func (s *CourseService) GetClassesByStudent(ctx context.Context, studentID string) ([]models.Class, error) {
	return s.repo.GetClassesByStudentID(ctx, studentID)
}

// GetClassesByLecturer gets classes for a lecturer
func (s *CourseService) GetClassesByLecturer(ctx context.Context, lecturerID string) ([]models.Class, error) {
	return s.repo.GetClassesByLecturerID(ctx, lecturerID)
}


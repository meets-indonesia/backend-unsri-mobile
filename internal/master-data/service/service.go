package service

import (
	"context"
	"time"

	"unsri-backend/internal/master-data/repository"
	apperrors "unsri-backend/internal/shared/errors"
	"unsri-backend/internal/shared/models"
)

// MasterDataService handles master data business logic
type MasterDataService struct {
	repo *repository.MasterDataRepository
}

// NewMasterDataService creates a new master data service
func NewMasterDataService(repo *repository.MasterDataRepository) *MasterDataService {
	return &MasterDataService{repo: repo}
}

// ========== Study Program Service Methods ==========

// CreateStudyProgramRequest represents create study program request
type CreateStudyProgramRequest struct {
	Code          string `json:"code" binding:"required"`
	Name          string `json:"name" binding:"required"`
	NameEn        string `json:"name_en,omitempty"`
	Faculty       string `json:"faculty,omitempty"`
	DegreeLevel   string `json:"degree_level,omitempty"`
	Accreditation string `json:"accreditation,omitempty"`
}

// CreateStudyProgram creates a new study program
func (s *MasterDataService) CreateStudyProgram(ctx context.Context, req CreateStudyProgramRequest) (*models.StudyProgram, error) {
	// Check if code already exists
	_, err := s.repo.GetStudyProgramByCode(ctx, req.Code)
	if err == nil {
		return nil, apperrors.NewConflictError("study program with code already exists")
	}

	studyProgram := &models.StudyProgram{
		Code:          req.Code,
		Name:          req.Name,
		NameEn:        req.NameEn,
		Faculty:       req.Faculty,
		DegreeLevel:   req.DegreeLevel,
		Accreditation: req.Accreditation,
		IsActive:      true,
	}

	if err := s.repo.CreateStudyProgram(ctx, studyProgram); err != nil {
		return nil, apperrors.NewInternalError("failed to create study program", err)
	}

	return studyProgram, nil
}

// GetStudyProgramByID gets a study program by ID
func (s *MasterDataService) GetStudyProgramByID(ctx context.Context, id string) (*models.StudyProgram, error) {
	studyProgram, err := s.repo.GetStudyProgramByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewNotFoundError("study program", id)
	}
	return studyProgram, nil
}

// GetStudyProgramsRequest represents get study programs request
type GetStudyProgramsRequest struct {
	Faculty  string `form:"faculty"`
	IsActive *bool  `form:"is_active"`
	Page     int    `form:"page,default=1"`
	PerPage  int    `form:"per_page,default=20"`
}

// GetStudyPrograms gets all study programs
func (s *MasterDataService) GetStudyPrograms(ctx context.Context, req GetStudyProgramsRequest) ([]models.StudyProgram, int64, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}

	var facultyPtr *string
	if req.Faculty != "" {
		facultyPtr = &req.Faculty
	}

	return s.repo.GetAllStudyPrograms(ctx, facultyPtr, req.IsActive, perPage, (page-1)*perPage)
}

// UpdateStudyProgramRequest represents update study program request
type UpdateStudyProgramRequest struct {
	Name          *string `json:"name,omitempty"`
	NameEn        *string `json:"name_en,omitempty"`
	Faculty       *string `json:"faculty,omitempty"`
	DegreeLevel   *string `json:"degree_level,omitempty"`
	Accreditation *string `json:"accreditation,omitempty"`
	IsActive      *bool   `json:"is_active,omitempty"`
}

// UpdateStudyProgram updates a study program
func (s *MasterDataService) UpdateStudyProgram(ctx context.Context, id string, req UpdateStudyProgramRequest) (*models.StudyProgram, error) {
	studyProgram, err := s.repo.GetStudyProgramByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewNotFoundError("study program", id)
	}

	if req.Name != nil {
		studyProgram.Name = *req.Name
	}
	if req.NameEn != nil {
		studyProgram.NameEn = *req.NameEn
	}
	if req.Faculty != nil {
		studyProgram.Faculty = *req.Faculty
	}
	if req.DegreeLevel != nil {
		studyProgram.DegreeLevel = *req.DegreeLevel
	}
	if req.Accreditation != nil {
		studyProgram.Accreditation = *req.Accreditation
	}
	if req.IsActive != nil {
		studyProgram.IsActive = *req.IsActive
	}

	if err := s.repo.UpdateStudyProgram(ctx, studyProgram); err != nil {
		return nil, apperrors.NewInternalError("failed to update study program", err)
	}

	return studyProgram, nil
}

// DeleteStudyProgram deletes a study program
func (s *MasterDataService) DeleteStudyProgram(ctx context.Context, id string) error {
	_, err := s.repo.GetStudyProgramByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFoundError("study program", id)
	}
	return s.repo.DeleteStudyProgram(ctx, id)
}

// ========== Academic Period Service Methods ==========

// CreateAcademicPeriodRequest represents create academic period request
type CreateAcademicPeriodRequest struct {
	Code              string  `json:"code" binding:"required"`
	Name              string  `json:"name" binding:"required"`
	AcademicYear      string  `json:"academic_year" binding:"required"`
	SemesterType      string  `json:"semester_type" binding:"required,oneof=GANJIL GENAP PENDEK"`
	StartDate         string  `json:"start_date" binding:"required"`
	EndDate           string  `json:"end_date" binding:"required"`
	RegistrationStart *string `json:"registration_start,omitempty"`
	RegistrationEnd   *string `json:"registration_end,omitempty"`
	IsActive          *bool   `json:"is_active,omitempty"`
}

// CreateAcademicPeriod creates a new academic period
func (s *MasterDataService) CreateAcademicPeriod(ctx context.Context, req CreateAcademicPeriodRequest) (*models.AcademicPeriod, error) {
	// Check if code already exists
	_, err := s.repo.GetAcademicPeriodByCode(ctx, req.Code)
	if err == nil {
		return nil, apperrors.NewConflictError("academic period with code already exists")
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, apperrors.NewValidationError("invalid start_date format, use YYYY-MM-DD")
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, apperrors.NewValidationError("invalid end_date format, use YYYY-MM-DD")
	}

	var registrationStart *time.Time
	if req.RegistrationStart != nil {
		regStart, err := time.Parse("2006-01-02", *req.RegistrationStart)
		if err != nil {
			return nil, apperrors.NewValidationError("invalid registration_start format, use YYYY-MM-DD")
		}
		registrationStart = &regStart
	}

	var registrationEnd *time.Time
	if req.RegistrationEnd != nil {
		regEnd, err := time.Parse("2006-01-02", *req.RegistrationEnd)
		if err != nil {
			return nil, apperrors.NewValidationError("invalid registration_end format, use YYYY-MM-DD")
		}
		registrationEnd = &regEnd
	}

	// If setting as active, deactivate other periods
	isActive := false
	if req.IsActive != nil && *req.IsActive {
		isActive = true
		// Deactivate other periods
		periods, _, err := s.repo.GetAllAcademicPeriods(ctx, nil, nil, nil, 1000, 0)
		if err == nil {
			for _, period := range periods {
				if period.IsActive {
					period.IsActive = false
					if err := s.repo.UpdateAcademicPeriod(ctx, &period); err != nil {
						return nil, apperrors.NewInternalError("failed to deactivate active period", err)
					}
				}
			}
		}
	}

	academicPeriod := &models.AcademicPeriod{
		Code:              req.Code,
		Name:              req.Name,
		AcademicYear:      req.AcademicYear,
		SemesterType:      req.SemesterType,
		StartDate:         startDate,
		EndDate:           endDate,
		RegistrationStart: registrationStart,
		RegistrationEnd:   registrationEnd,
		IsActive:          isActive,
	}

	if err := s.repo.CreateAcademicPeriod(ctx, academicPeriod); err != nil {
		return nil, apperrors.NewInternalError("failed to create academic period", err)
	}

	return academicPeriod, nil
}

// GetAcademicPeriodByID gets an academic period by ID
func (s *MasterDataService) GetAcademicPeriodByID(ctx context.Context, id string) (*models.AcademicPeriod, error) {
	academicPeriod, err := s.repo.GetAcademicPeriodByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewNotFoundError("academic period", id)
	}
	return academicPeriod, nil
}

// GetActiveAcademicPeriod gets the active academic period
func (s *MasterDataService) GetActiveAcademicPeriod(ctx context.Context) (*models.AcademicPeriod, error) {
	academicPeriod, err := s.repo.GetActiveAcademicPeriod(ctx)
	if err != nil {
		return nil, apperrors.NewNotFoundError("active academic period", "")
	}
	return academicPeriod, nil
}

// GetAcademicPeriodsRequest represents get academic periods request
type GetAcademicPeriodsRequest struct {
	AcademicYear string `form:"academic_year"`
	SemesterType string `form:"semester_type"`
	IsActive     *bool  `form:"is_active"`
	Page         int    `form:"page,default=1"`
	PerPage      int    `form:"per_page,default=20"`
}

// GetAcademicPeriods gets all academic periods
func (s *MasterDataService) GetAcademicPeriods(ctx context.Context, req GetAcademicPeriodsRequest) ([]models.AcademicPeriod, int64, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}

	var academicYearPtr, semesterTypePtr *string
	if req.AcademicYear != "" {
		academicYearPtr = &req.AcademicYear
	}
	if req.SemesterType != "" {
		semesterTypePtr = &req.SemesterType
	}

	return s.repo.GetAllAcademicPeriods(ctx, academicYearPtr, semesterTypePtr, req.IsActive, perPage, (page-1)*perPage)
}

// UpdateAcademicPeriodRequest represents update academic period request
type UpdateAcademicPeriodRequest struct {
	Name              *string `json:"name,omitempty"`
	AcademicYear      *string `json:"academic_year,omitempty"`
	SemesterType      *string `json:"semester_type,omitempty"`
	StartDate         *string `json:"start_date,omitempty"`
	EndDate           *string `json:"end_date,omitempty"`
	RegistrationStart *string `json:"registration_start,omitempty"`
	RegistrationEnd   *string `json:"registration_end,omitempty"`
	IsActive          *bool   `json:"is_active,omitempty"`
}

// UpdateAcademicPeriod updates an academic period
func (s *MasterDataService) UpdateAcademicPeriod(ctx context.Context, id string, req UpdateAcademicPeriodRequest) (*models.AcademicPeriod, error) {
	academicPeriod, err := s.repo.GetAcademicPeriodByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewNotFoundError("academic period", id)
	}

	if req.Name != nil {
		academicPeriod.Name = *req.Name
	}
	if req.AcademicYear != nil {
		academicPeriod.AcademicYear = *req.AcademicYear
	}
	if req.SemesterType != nil {
		academicPeriod.SemesterType = *req.SemesterType
	}
	if req.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			return nil, apperrors.NewValidationError("invalid start_date format, use YYYY-MM-DD")
		}
		academicPeriod.StartDate = startDate
	}
	if req.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return nil, apperrors.NewValidationError("invalid end_date format, use YYYY-MM-DD")
		}
		academicPeriod.EndDate = endDate
	}
	if req.RegistrationStart != nil {
		regStart, err := time.Parse("2006-01-02", *req.RegistrationStart)
		if err != nil {
			return nil, apperrors.NewValidationError("invalid registration_start format, use YYYY-MM-DD")
		}
		academicPeriod.RegistrationStart = &regStart
	}
	if req.RegistrationEnd != nil {
		regEnd, err := time.Parse("2006-01-02", *req.RegistrationEnd)
		if err != nil {
			return nil, apperrors.NewValidationError("invalid registration_end format, use YYYY-MM-DD")
		}
		academicPeriod.RegistrationEnd = &regEnd
	}
	if req.IsActive != nil && *req.IsActive {
		// If setting as active, deactivate other periods
		periods, _, err := s.repo.GetAllAcademicPeriods(ctx, nil, nil, nil, 1000, 0)
		if err == nil {
			for _, period := range periods {
				if period.IsActive && period.ID != id {
					period.IsActive = false
					if err := s.repo.UpdateAcademicPeriod(ctx, &period); err != nil {
						return nil, apperrors.NewInternalError("failed to deactivate active period", err)
					}
				}
			}
		}
		academicPeriod.IsActive = true
	} else if req.IsActive != nil {
		academicPeriod.IsActive = false
	}

	if err := s.repo.UpdateAcademicPeriod(ctx, academicPeriod); err != nil {
		return nil, apperrors.NewInternalError("failed to update academic period", err)
	}

	return academicPeriod, nil
}

// DeleteAcademicPeriod deletes an academic period
func (s *MasterDataService) DeleteAcademicPeriod(ctx context.Context, id string) error {
	_, err := s.repo.GetAcademicPeriodByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFoundError("academic period", id)
	}
	return s.repo.DeleteAcademicPeriod(ctx, id)
}

// ========== Room Service Methods ==========

// CreateRoomRequest represents create room request
type CreateRoomRequest struct {
	Code       string `json:"code" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Building   string `json:"building,omitempty"`
	Floor      *int   `json:"floor,omitempty"`
	Capacity   *int   `json:"capacity,omitempty"`
	RoomType   string `json:"room_type,omitempty"`
	Facilities string `json:"facilities,omitempty"`
}

// CreateRoom creates a new room
func (s *MasterDataService) CreateRoom(ctx context.Context, req CreateRoomRequest) (*models.Room, error) {
	// Check if code already exists
	_, err := s.repo.GetRoomByCode(ctx, req.Code)
	if err == nil {
		return nil, apperrors.NewConflictError("room with code already exists")
	}

	room := &models.Room{
		Code:       req.Code,
		Name:       req.Name,
		Building:   req.Building,
		Floor:      req.Floor,
		Capacity:   req.Capacity,
		RoomType:   req.RoomType,
		Facilities: req.Facilities,
		IsActive:   true,
	}

	if err := s.repo.CreateRoom(ctx, room); err != nil {
		return nil, apperrors.NewInternalError("failed to create room", err)
	}

	return room, nil
}

// GetRoomByID gets a room by ID
func (s *MasterDataService) GetRoomByID(ctx context.Context, id string) (*models.Room, error) {
	room, err := s.repo.GetRoomByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewNotFoundError("room", id)
	}
	return room, nil
}

// GetRoomsRequest represents get rooms request
type GetRoomsRequest struct {
	Building string `form:"building"`
	RoomType string `form:"room_type"`
	IsActive *bool  `form:"is_active"`
	Page     int    `form:"page,default=1"`
	PerPage  int    `form:"per_page,default=20"`
}

// GetRooms gets all rooms
func (s *MasterDataService) GetRooms(ctx context.Context, req GetRoomsRequest) ([]models.Room, int64, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}

	var buildingPtr, roomTypePtr *string
	if req.Building != "" {
		buildingPtr = &req.Building
	}
	if req.RoomType != "" {
		roomTypePtr = &req.RoomType
	}

	return s.repo.GetAllRooms(ctx, buildingPtr, roomTypePtr, req.IsActive, perPage, (page-1)*perPage)
}

// UpdateRoomRequest represents update room request
type UpdateRoomRequest struct {
	Name       *string `json:"name,omitempty"`
	Building   *string `json:"building,omitempty"`
	Floor      *int    `json:"floor,omitempty"`
	Capacity   *int    `json:"capacity,omitempty"`
	RoomType   *string `json:"room_type,omitempty"`
	Facilities *string `json:"facilities,omitempty"`
	IsActive   *bool   `json:"is_active,omitempty"`
}

// UpdateRoom updates a room
func (s *MasterDataService) UpdateRoom(ctx context.Context, id string, req UpdateRoomRequest) (*models.Room, error) {
	room, err := s.repo.GetRoomByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewNotFoundError("room", id)
	}

	if req.Name != nil {
		room.Name = *req.Name
	}
	if req.Building != nil {
		room.Building = *req.Building
	}
	if req.Floor != nil {
		room.Floor = req.Floor
	}
	if req.Capacity != nil {
		room.Capacity = req.Capacity
	}
	if req.RoomType != nil {
		room.RoomType = *req.RoomType
	}
	if req.Facilities != nil {
		room.Facilities = *req.Facilities
	}
	if req.IsActive != nil {
		room.IsActive = *req.IsActive
	}

	if err := s.repo.UpdateRoom(ctx, room); err != nil {
		return nil, apperrors.NewInternalError("failed to update room", err)
	}

	return room, nil
}

// DeleteRoom deletes a room
func (s *MasterDataService) DeleteRoom(ctx context.Context, id string) error {
	_, err := s.repo.GetRoomByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFoundError("room", id)
	}
	return s.repo.DeleteRoom(ctx, id)
}

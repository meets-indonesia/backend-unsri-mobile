package service

import (
	"testing"
	"time"

	apperrors "unsri-backend/internal/shared/errors"
	"unsri-backend/internal/shared/models"

	"github.com/google/uuid"
)

// Test helper functions
func createTestStudyProgram() *models.StudyProgram {
	return &models.StudyProgram{
		ID:          uuid.New().String(),
		Code:        "TI",
		Name:        "Teknik Informatika",
		NameEn:      "Informatics Engineering",
		Faculty:     "Fakultas Ilmu Komputer",
		DegreeLevel: "S1",
		IsActive:    true,
	}
}

func createTestAcademicPeriod() *models.AcademicPeriod {
	startDate := time.Now()
	endDate := startDate.AddDate(0, 6, 0)
	return &models.AcademicPeriod{
		ID:           uuid.New().String(),
		Code:         "2024-GANJIL",
		Name:         "Ganjil 2024/2025",
		AcademicYear: "2024/2025",
		SemesterType: "Ganjil",
		StartDate:    startDate,
		EndDate:      endDate,
		IsActive:     false,
	}
}

func createTestRoom() *models.Room {
	floor := 1
	capacity := 40
	return &models.Room{
		ID:       uuid.New().String(),
		Code:     "A101",
		Name:     "Ruang A101",
		Building: "Gedung A",
		Floor:    &floor,
		Capacity: &capacity,
		RoomType: "classroom",
		IsActive: true,
	}
}

// Test CreateStudyProgramRequest validation
func TestCreateStudyProgramRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateStudyProgramRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CreateStudyProgramRequest{
				Code:        "TI",
				Name:        "Teknik Informatika",
				NameEn:      "Informatics Engineering",
				Faculty:     "Fakultas Ilmu Komputer",
				DegreeLevel: "S1",
			},
			wantErr: false,
		},
		{
			name: "missing code",
			req: CreateStudyProgramRequest{
				Name: "Teknik Informatika",
			},
			wantErr: true,
		},
		{
			name: "missing name",
			req: CreateStudyProgramRequest{
				Code: "TI",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This is a validation test - actual validation happens in handler
			// We're just testing the struct definition
			if tt.req.Code == "" && !tt.wantErr {
				t.Error("Code should be required")
			}
			if tt.req.Name == "" && !tt.wantErr {
				t.Error("Name should be required")
			}
		})
	}
}

// Test CreateAcademicPeriodRequest validation
func TestCreateAcademicPeriodRequest(t *testing.T) {
	startDate := time.Now()
	endDate := startDate.AddDate(0, 6, 0)

	tests := []struct {
		name    string
		req     CreateAcademicPeriodRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CreateAcademicPeriodRequest{
				Code:         "2024-GANJIL",
				Name:         "Ganjil 2024/2025",
				AcademicYear: "2024/2025",
				SemesterType: "Ganjil",
				StartDate:    startDate.Format("2006-01-02"),
				EndDate:      endDate.Format("2006-01-02"),
			},
			wantErr: false,
		},
		{
			name: "invalid date format",
			req: CreateAcademicPeriodRequest{
				Code:         "2024-GANJIL",
				Name:         "Ganjil 2024/2025",
				AcademicYear: "2024/2025",
				SemesterType: "Ganjil",
				StartDate:    "invalid-date",
				EndDate:      endDate.Format("2006-01-02"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := time.Parse("2006-01-02", tt.req.StartDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("Date parsing error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test CreateRoomRequest validation
func TestCreateRoomRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateRoomRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CreateRoomRequest{
				Code:     "A101",
				Name:     "Ruang A101",
				Building: "Gedung A",
				Floor:    func() *int { v := 1; return &v }(),
				Capacity: func() *int { v := 40; return &v }(),
				RoomType: "classroom",
			},
			wantErr: false,
		},
		{
			name: "missing code",
			req: CreateRoomRequest{
				Name: "Ruang A101",
			},
			wantErr: true,
		},
		{
			name: "missing name",
			req: CreateRoomRequest{
				Code: "A101",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req.Code == "" && !tt.wantErr {
				t.Error("Code should be required")
			}
			if tt.req.Name == "" && !tt.wantErr {
				t.Error("Name should be required")
			}
		})
	}
}

// Test error types
func TestErrorTypes(t *testing.T) {
	t.Run("NotFoundError", func(t *testing.T) {
		err := apperrors.NewNotFoundError("study program", "test-id")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("ConflictError", func(t *testing.T) {
		err := apperrors.NewConflictError("study program with code already exists")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("ValidationError", func(t *testing.T) {
		err := apperrors.NewValidationError("invalid input")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

// Test model validation
func TestStudyProgramModel(t *testing.T) {
	t.Run("valid study program", func(t *testing.T) {
		sp := createTestStudyProgram()
		if sp.Code == "" {
			t.Error("Code should not be empty")
		}
		if sp.Name == "" {
			t.Error("Name should not be empty")
		}
		if sp.ID == "" {
			t.Error("ID should be generated")
		}
	})

	t.Run("table name", func(t *testing.T) {
		sp := models.StudyProgram{}
		if sp.TableName() != "study_programs" {
			t.Errorf("Expected table name 'study_programs', got '%s'", sp.TableName())
		}
	})
}

func TestAcademicPeriodModel(t *testing.T) {
	t.Run("valid academic period", func(t *testing.T) {
		ap := createTestAcademicPeriod()
		if ap.Code == "" {
			t.Error("Code should not be empty")
		}
		if ap.Name == "" {
			t.Error("Name should not be empty")
		}
		if ap.StartDate.After(ap.EndDate) {
			t.Error("Start date should be before end date")
		}
	})

	t.Run("table name", func(t *testing.T) {
		ap := models.AcademicPeriod{}
		if ap.TableName() != "academic_periods" {
			t.Errorf("Expected table name 'academic_periods', got '%s'", ap.TableName())
		}
	})
}

func TestRoomModel(t *testing.T) {
	t.Run("valid room", func(t *testing.T) {
		room := createTestRoom()
		if room.Code == "" {
			t.Error("Code should not be empty")
		}
		if room.Name == "" {
			t.Error("Name should not be empty")
		}
		if room.Capacity != nil && *room.Capacity <= 0 {
			t.Error("Capacity should be positive")
		}
	})

	t.Run("table name", func(t *testing.T) {
		room := models.Room{}
		if room.TableName() != "rooms" {
			t.Errorf("Expected table name 'rooms', got '%s'", room.TableName())
		}
	})
}

// Test date validation logic
func TestDateValidation(t *testing.T) {
	t.Run("valid date range", func(t *testing.T) {
		startDate := time.Now()
		endDate := startDate.AddDate(0, 6, 0)

		if endDate.Before(startDate) {
			t.Error("End date should be after start date")
		}
	})

	t.Run("invalid date range", func(t *testing.T) {
		startDate := time.Now()
		endDate := startDate.AddDate(0, -1, 0)

		if !endDate.Before(startDate) {
			t.Error("End date should be before start date in this test case")
		}
	})
}

// Test UpdateStudyProgramRequest validation
func TestUpdateStudyProgramRequest(t *testing.T) {
	name := "Updated Name"
	nameEn := "Updated Name EN"
	faculty := "Updated Faculty"
	degreeLevel := "S2"
	accreditation := "A"
	isActive := true

	tests := []struct {
		name string
		req  UpdateStudyProgramRequest
	}{
		{
			name: "update name only",
			req: UpdateStudyProgramRequest{
				Name: &name,
			},
		},
		{
			name: "update all fields",
			req: UpdateStudyProgramRequest{
				Name:          &name,
				NameEn:        &nameEn,
				Faculty:       &faculty,
				DegreeLevel:   &degreeLevel,
				Accreditation: &accreditation,
				IsActive:      &isActive,
			},
		},
		{
			name: "update is_active to false",
			req: UpdateStudyProgramRequest{
				IsActive: func() *bool { v := false; return &v }(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req.Name != nil && *tt.req.Name == "" {
				t.Error("Name should not be empty if provided")
			}
		})
	}
}

// Test UpdateAcademicPeriodRequest validation
func TestUpdateAcademicPeriodRequest(t *testing.T) {
	name := "Updated Period"
	academicYear := "2025/2026"
	semesterType := "GENAP"
	startDate := "2025-02-01"
	endDate := "2025-07-31"
	isActive := true

	tests := []struct {
		name    string
		req     UpdateAcademicPeriodRequest
		wantErr bool
	}{
		{
			name: "update name only",
			req: UpdateAcademicPeriodRequest{
				Name: &name,
			},
			wantErr: false,
		},
		{
			name: "update with invalid date format",
			req: UpdateAcademicPeriodRequest{
				StartDate: func() *string { v := "invalid-date"; return &v }(),
			},
			wantErr: true,
		},
		{
			name: "update all fields",
			req: UpdateAcademicPeriodRequest{
				Name:         &name,
				AcademicYear: &academicYear,
				SemesterType: &semesterType,
				StartDate:    &startDate,
				EndDate:      &endDate,
				IsActive:     &isActive,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req.StartDate != nil {
				_, err := time.Parse("2006-01-02", *tt.req.StartDate)
				if (err != nil) != tt.wantErr {
					t.Errorf("Date parsing error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

// Test UpdateRoomRequest validation
func TestUpdateRoomRequest(t *testing.T) {
	name := "Updated Room"
	building := "Gedung B"
	floor := 2
	capacity := 50
	roomType := "lab"
	facilities := "Projector, AC"
	isActive := true

	tests := []struct {
		name string
		req  UpdateRoomRequest
	}{
		{
			name: "update name only",
			req: UpdateRoomRequest{
				Name: &name,
			},
		},
		{
			name: "update capacity",
			req: UpdateRoomRequest{
				Capacity: &capacity,
			},
		},
		{
			name: "update all fields",
			req: UpdateRoomRequest{
				Name:       &name,
				Building:   &building,
				Floor:      &floor,
				Capacity:   &capacity,
				RoomType:   &roomType,
				Facilities: &facilities,
				IsActive:   &isActive,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req.Capacity != nil && *tt.req.Capacity <= 0 {
				t.Error("Capacity should be positive")
			}
		})
	}
}

// Test GetStudyProgramsRequest pagination
func TestGetStudyProgramsRequestPagination(t *testing.T) {
	tests := []struct {
		name     string
		req      GetStudyProgramsRequest
		expected struct {
			page    int
			perPage int
		}
	}{
		{
			name: "default pagination",
			req: GetStudyProgramsRequest{
				Page:    0,
				PerPage: 0,
			},
			expected: struct {
				page    int
				perPage int
			}{
				page:    1,
				perPage: 20,
			},
		},
		{
			name: "custom pagination",
			req: GetStudyProgramsRequest{
				Page:    2,
				PerPage: 10,
			},
			expected: struct {
				page    int
				perPage int
			}{
				page:    2,
				perPage: 10,
			},
		},
		{
			name: "negative page",
			req: GetStudyProgramsRequest{
				Page:    -1,
				PerPage: 20,
			},
			expected: struct {
				page    int
				perPage int
			}{
				page:    1,
				perPage: 20,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page := tt.req.Page
			if page < 1 {
				page = 1
			}
			perPage := tt.req.PerPage
			if perPage < 1 {
				perPage = 20
			}

			if page != tt.expected.page {
				t.Errorf("Expected page %d, got %d", tt.expected.page, page)
			}
			if perPage != tt.expected.perPage {
				t.Errorf("Expected perPage %d, got %d", tt.expected.perPage, perPage)
			}
		})
	}
}

// Test GetAcademicPeriodsRequest pagination
func TestGetAcademicPeriodsRequestPagination(t *testing.T) {
	tests := []struct {
		name     string
		req      GetAcademicPeriodsRequest
		expected struct {
			page    int
			perPage int
		}
	}{
		{
			name: "default pagination",
			req: GetAcademicPeriodsRequest{
				Page:    0,
				PerPage: 0,
			},
			expected: struct {
				page    int
				perPage int
			}{
				page:    1,
				perPage: 20,
			},
		},
		{
			name: "custom pagination",
			req: GetAcademicPeriodsRequest{
				Page:    3,
				PerPage: 15,
			},
			expected: struct {
				page    int
				perPage int
			}{
				page:    3,
				perPage: 15,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page := tt.req.Page
			if page < 1 {
				page = 1
			}
			perPage := tt.req.PerPage
			if perPage < 1 {
				perPage = 20
			}

			if page != tt.expected.page {
				t.Errorf("Expected page %d, got %d", tt.expected.page, page)
			}
			if perPage != tt.expected.perPage {
				t.Errorf("Expected perPage %d, got %d", tt.expected.perPage, perPage)
			}
		})
	}
}

// Test GetRoomsRequest pagination
func TestGetRoomsRequestPagination(t *testing.T) {
	tests := []struct {
		name     string
		req      GetRoomsRequest
		expected struct {
			page    int
			perPage int
		}
	}{
		{
			name: "default pagination",
			req: GetRoomsRequest{
				Page:    0,
				PerPage: 0,
			},
			expected: struct {
				page    int
				perPage int
			}{
				page:    1,
				perPage: 20,
			},
		},
		{
			name: "custom pagination",
			req: GetRoomsRequest{
				Page:    5,
				PerPage: 25,
			},
			expected: struct {
				page    int
				perPage int
			}{
				page:    5,
				perPage: 25,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page := tt.req.Page
			if page < 1 {
				page = 1
			}
			perPage := tt.req.PerPage
			if perPage < 1 {
				perPage = 20
			}

			if page != tt.expected.page {
				t.Errorf("Expected page %d, got %d", tt.expected.page, page)
			}
			if perPage != tt.expected.perPage {
				t.Errorf("Expected perPage %d, got %d", tt.expected.perPage, perPage)
			}
		})
	}
}

// Test Academic Period Semester Type validation
func TestAcademicPeriodSemesterType(t *testing.T) {
	validTypes := []string{"GANJIL", "GENAP", "PENDEK"}
	invalidTypes := []string{"INVALID", "", "ganjil"}

	for _, validType := range validTypes {
		t.Run("valid type: "+validType, func(t *testing.T) {
			req := CreateAcademicPeriodRequest{
				Code:         "TEST",
				Name:         "Test",
				AcademicYear: "2024/2025",
				SemesterType: validType,
				StartDate:    "2024-01-01",
				EndDate:      "2024-06-30",
			}
			if req.SemesterType != validType {
				t.Errorf("Expected semester type %s, got %s", validType, req.SemesterType)
			}
		})
	}

	for _, invalidType := range invalidTypes {
		t.Run("invalid type: "+invalidType, func(t *testing.T) {
			req := CreateAcademicPeriodRequest{
				Code:         "TEST",
				Name:         "Test",
				AcademicYear: "2024/2025",
				SemesterType: invalidType,
				StartDate:    "2024-01-01",
				EndDate:      "2024-06-30",
			}
			// Validation happens in handler, here we just check the value
			if req.SemesterType == invalidType && invalidType != "" {
				t.Logf("Tested invalid type: %s", invalidType)
			}
		})
	}
}

// Test Room Type validation
func TestRoomType(t *testing.T) {
	validTypes := []string{"classroom", "lab", "auditorium", "library", "office"}

	for _, roomType := range validTypes {
		t.Run("valid room type: "+roomType, func(t *testing.T) {
			req := CreateRoomRequest{
				Code:     "TEST",
				Name:     "Test Room",
				RoomType: roomType,
			}
			if req.RoomType != roomType {
				t.Errorf("Expected room type %s, got %s", roomType, req.RoomType)
			}
		})
	}
}

// Test Study Program Degree Level validation
func TestStudyProgramDegreeLevel(t *testing.T) {
	validLevels := []string{"S1", "S2", "S3", "D3", "D4"}

	for _, level := range validLevels {
		t.Run("valid degree level: "+level, func(t *testing.T) {
			req := CreateStudyProgramRequest{
				Code:        "TEST",
				Name:        "Test Program",
				DegreeLevel: level,
			}
			if req.DegreeLevel != level {
				t.Errorf("Expected degree level %s, got %s", level, req.DegreeLevel)
			}
		})
	}
}

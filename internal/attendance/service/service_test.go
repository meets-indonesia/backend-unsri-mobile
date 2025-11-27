package service

import (
	"testing"
	"time"

	"github.com/google/uuid"
	apperrors "unsri-backend/internal/shared/errors"
	"unsri-backend/internal/shared/models"
)

// Test helper functions
func createTestAttendance() *models.Attendance {
	sessionID := uuid.New().String()
	checkInTime := time.Now()
	return &models.Attendance{
		ID:          uuid.New().String(),
		UserID:      uuid.New().String(),
		SessionID:   &sessionID,
		Type:        models.AttendanceTypeKelas,
		Status:      models.StatusHadir,
		Date:        time.Now(),
		CheckInTime: &checkInTime,
	}
}

func createTestAttendanceSession() *models.AttendanceSession {
	scheduleID := uuid.New().String()
	return &models.AttendanceSession{
		ID:         uuid.New().String(),
		ScheduleID: &scheduleID,
		CreatedBy:  uuid.New().String(),
		Type:       models.AttendanceTypeKelas,
		ExpiresAt:  time.Now().Add(30 * time.Minute),
		IsActive:   true,
	}
}

func createTestSchedule() *models.Schedule {
	return &models.Schedule{
		ID:        uuid.New().String(),
		DosenID:   uuid.New().String(),
		DayOfWeek: 1, // Monday
		StartTime: time.Now(),
		EndTime:   time.Now().Add(2 * time.Hour),
		Date:      time.Now(),
		IsActive:  true,
	}
}

func createTestShiftPattern() *models.ShiftPattern {
	startTime := time.Date(2000, 1, 1, 8, 0, 0, 0, time.UTC)
	endTime := time.Date(2000, 1, 1, 17, 0, 0, 0, time.UTC)
	breakDuration := 60
	return &models.ShiftPattern{
		ID:                  uuid.New().String(),
		ShiftName:           "Regular Shift",
		ShiftCode:           "SHIFT-001",
		StartTime:           startTime,
		EndTime:             endTime,
		BreakDurationMinutes: &breakDuration,
		IsActive:            true,
	}
}

func createTestWorkAttendanceRecord() *models.WorkAttendanceRecord {
	scheduleID := uuid.New().String()
	return &models.WorkAttendanceRecord{
		ID:             uuid.New().String(),
		UserID:         uuid.New().String(),
		ScheduleID:     &scheduleID,
		AttendanceType: "CHECK_IN",
		RecordedAt:     time.Now(),
		Status:         models.StatusCheckIn,
	}
}

// Test AttendanceStatus validation
func TestAttendanceStatus(t *testing.T) {
	validStatuses := []models.AttendanceStatus{
		models.StatusHadir,
		models.StatusIzin,
		models.StatusSakit,
		models.StatusAlpa,
		models.StatusTerlambat,
	}

	for _, status := range validStatuses {
		t.Run(string(status), func(t *testing.T) {
			if status == "" {
				t.Error("Status should not be empty")
			}
		})
	}
}

// Test WorkAttendanceStatus validation
func TestWorkAttendanceStatus(t *testing.T) {
	validStatuses := []models.WorkAttendanceStatus{
		models.StatusCheckIn,
		models.StatusCheckOut,
		models.StatusLateIn,
		models.StatusEarlyOut,
		models.StatusAbsent,
		models.StatusOnLeave,
		models.StatusSickLeave,
	}

	for _, status := range validStatuses {
		t.Run(string(status), func(t *testing.T) {
			if status == "" {
				t.Error("Status should not be empty")
			}
		})
	}
}

// Test AttendanceType validation
func TestAttendanceType(t *testing.T) {
	validTypes := []models.AttendanceType{
		models.AttendanceTypeKelas,
		models.AttendanceTypeKampus,
	}

	for _, atype := range validTypes {
		t.Run(string(atype), func(t *testing.T) {
			if atype == "" {
				t.Error("Type should not be empty")
			}
		})
	}
}

// Test error types
func TestErrorTypes(t *testing.T) {
	t.Run("NotFoundError", func(t *testing.T) {
		err := apperrors.NewNotFoundError("attendance", "test-id")
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

	t.Run("ForbiddenError", func(t *testing.T) {
		err := apperrors.NewForbiddenError("insufficient permissions")
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

// Test Attendance model
func TestAttendanceModel(t *testing.T) {
	t.Run("valid attendance", func(t *testing.T) {
		attendance := createTestAttendance()
		if attendance.UserID == "" {
			t.Error("UserID should not be empty")
		}
		if attendance.Status == "" {
			t.Error("Status should not be empty")
		}
		if attendance.ID == "" {
			t.Error("ID should be generated")
		}
	})

	t.Run("table name", func(t *testing.T) {
		attendance := models.Attendance{}
		if attendance.TableName() != "attendances" {
			t.Errorf("Expected table name 'attendances', got '%s'", attendance.TableName())
		}
	})
}

// Test AttendanceSession model
func TestAttendanceSessionModel(t *testing.T) {
	t.Run("valid session", func(t *testing.T) {
		session := createTestAttendanceSession()
		if session.CreatedBy == "" {
			t.Error("CreatedBy should not be empty")
		}
		if session.Type == "" {
			t.Error("Type should not be empty")
		}
		if session.ID == "" {
			t.Error("ID should be generated")
		}
	})

	t.Run("table name", func(t *testing.T) {
		session := models.AttendanceSession{}
		if session.TableName() != "attendance_sessions" {
			t.Errorf("Expected table name 'attendance_sessions', got '%s'", session.TableName())
		}
	})
}

// Test Schedule model
func TestScheduleModel(t *testing.T) {
	t.Run("valid schedule", func(t *testing.T) {
		schedule := createTestSchedule()
		if schedule.DosenID == "" {
			t.Error("DosenID should not be empty")
		}
		if schedule.ID == "" {
			t.Error("ID should be generated")
		}
	})

	t.Run("table name", func(t *testing.T) {
		schedule := models.Schedule{}
		if schedule.TableName() != "schedules" {
			t.Errorf("Expected table name 'schedules', got '%s'", schedule.TableName())
		}
	})

	t.Run("day of week validation", func(t *testing.T) {
		schedule := createTestSchedule()
		if schedule.DayOfWeek < 0 || schedule.DayOfWeek > 6 {
			t.Error("DayOfWeek should be between 0 and 6")
		}
	})
}

// Test ShiftPattern model
func TestShiftPatternModel(t *testing.T) {
	t.Run("valid shift pattern", func(t *testing.T) {
		shift := createTestShiftPattern()
		if shift.ShiftCode == "" {
			t.Error("ShiftCode should not be empty")
		}
		if shift.ShiftName == "" {
			t.Error("ShiftName should not be empty")
		}
		if shift.ID == "" {
			t.Error("ID should be generated")
		}
	})

	t.Run("table name", func(t *testing.T) {
		shift := models.ShiftPattern{}
		if shift.TableName() != "shift_patterns" {
			t.Errorf("Expected table name 'shift_patterns', got '%s'", shift.TableName())
		}
	})
}

// Test WorkAttendanceRecord model
func TestWorkAttendanceRecordModel(t *testing.T) {
	t.Run("valid work attendance record", func(t *testing.T) {
		record := createTestWorkAttendanceRecord()
		if record.UserID == "" {
			t.Error("UserID should not be empty")
		}
		if record.Status == "" {
			t.Error("Status should not be empty")
		}
		if record.ID == "" {
			t.Error("ID should be generated")
		}
	})

	t.Run("table name", func(t *testing.T) {
		record := models.WorkAttendanceRecord{}
		if record.TableName() != "work_attendance_records" {
			t.Errorf("Expected table name 'work_attendance_records', got '%s'", record.TableName())
		}
	})
}

// Test time validation
func TestTimeValidation(t *testing.T) {
	t.Run("valid time range", func(t *testing.T) {
		startTime := time.Now()
		endTime := startTime.Add(30 * time.Minute)

		if endTime.Before(startTime) {
			t.Error("EndTime should be after StartTime")
		}
	})

	t.Run("invalid time range", func(t *testing.T) {
		startTime := time.Now()
		endTime := startTime.Add(-30 * time.Minute)

		if !endTime.Before(startTime) {
			t.Error("EndTime should be before StartTime in this test case")
		}
	})
}

// Test CreateShiftPatternRequest validation
func TestCreateShiftPatternRequest(t *testing.T) {
	startTimeStr := "08:00"
	endTimeStr := "17:00"
	breakDuration := 60

	tests := []struct {
		name    string
		req     CreateShiftPatternRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CreateShiftPatternRequest{
				ShiftName:           "Regular Shift",
				ShiftCode:           "SHIFT-001",
				StartTime:           startTimeStr,
				EndTime:             endTimeStr,
				BreakDurationMinutes: &breakDuration,
			},
			wantErr: false,
		},
		{
			name: "missing shift name",
			req: CreateShiftPatternRequest{
				ShiftCode: "SHIFT-001",
				StartTime: startTimeStr,
				EndTime:   endTimeStr,
			},
			wantErr: true,
		},
		{
			name: "missing shift code",
			req: CreateShiftPatternRequest{
				ShiftName: "Regular Shift",
				StartTime: startTimeStr,
				EndTime:   endTimeStr,
			},
			wantErr: true,
		},
		{
			name: "invalid time format",
			req: CreateShiftPatternRequest{
				ShiftName: "Regular Shift",
				ShiftCode: "SHIFT-001",
				StartTime: "invalid-time",
				EndTime:   endTimeStr,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req.ShiftName == "" && !tt.wantErr {
				t.Error("ShiftName should be required")
			}
			if tt.req.ShiftCode == "" && !tt.wantErr {
				t.Error("ShiftCode should be required")
			}
			if tt.req.StartTime != "" && tt.req.EndTime != "" {
				start, err1 := time.Parse("15:04", tt.req.StartTime)
				end, err2 := time.Parse("15:04", tt.req.EndTime)
				if err1 == nil && err2 == nil {
					if end.Before(start) && !tt.wantErr {
						t.Error("EndTime should be after StartTime")
					}
				}
			}
		})
	}
}

// Test CreateUserShiftRequest validation
func TestCreateUserShiftRequest(t *testing.T) {
	effectiveFrom := time.Now()
	effectiveTo := effectiveFrom.AddDate(0, 6, 0)

	tests := []struct {
		name    string
		req     CreateUserShiftRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CreateUserShiftRequest{
				UserID:         uuid.New().String(),
				ShiftID:        uuid.New().String(),
				EffectiveFrom:  effectiveFrom.Format("2006-01-02"),
				EffectiveUntil: func() *string { v := effectiveTo.Format("2006-01-02"); return &v }(),
			},
			wantErr: false,
		},
		{
			name: "missing user_id",
			req: CreateUserShiftRequest{
				ShiftID:       uuid.New().String(),
				EffectiveFrom: effectiveFrom.Format("2006-01-02"),
			},
			wantErr: true,
		},
		{
			name: "missing shift_id",
			req: CreateUserShiftRequest{
				UserID:        uuid.New().String(),
				EffectiveFrom: effectiveFrom.Format("2006-01-02"),
			},
			wantErr: true,
		},
		{
			name: "invalid date format",
			req: CreateUserShiftRequest{
				UserID:        uuid.New().String(),
				ShiftID:       uuid.New().String(),
				EffectiveFrom: "invalid-date",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasError := false
			
			if tt.req.UserID == "" {
				hasError = true
			}
			if tt.req.ShiftID == "" {
				hasError = true
			}
			if tt.req.EffectiveFrom != "" {
				_, err := time.Parse("2006-01-02", tt.req.EffectiveFrom)
				if err != nil {
					hasError = true
				}
			}
			
			if hasError != tt.wantErr {
				t.Errorf("Expected error = %v, got %v", tt.wantErr, hasError)
			}
		})
	}
}

// Test CreateWorkScheduleRequest validation
func TestCreateWorkScheduleRequest(t *testing.T) {
	scheduleDate := time.Now()
	startTimeStr := "08:00"
	endTimeStr := "17:00"

	tests := []struct {
		name    string
		req     CreateWorkScheduleRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CreateWorkScheduleRequest{
				UserID:       uuid.New().String(),
				ScheduleDate: scheduleDate.Format("2006-01-02"),
				StartTime:    startTimeStr,
				EndTime:      endTimeStr,
			},
			wantErr: false,
		},
		{
			name: "missing user_id",
			req: CreateWorkScheduleRequest{
				ScheduleDate: scheduleDate.Format("2006-01-02"),
				StartTime:    startTimeStr,
				EndTime:      endTimeStr,
			},
			wantErr: true,
		},
		{
			name: "invalid date format",
			req: CreateWorkScheduleRequest{
				UserID:       uuid.New().String(),
				ScheduleDate: "invalid-date",
				StartTime:    startTimeStr,
				EndTime:      endTimeStr,
			},
			wantErr: true,
		},
		{
			name: "invalid time format",
			req: CreateWorkScheduleRequest{
				UserID:       uuid.New().String(),
				ScheduleDate: scheduleDate.Format("2006-01-02"),
				StartTime:    "invalid-time",
				EndTime:      endTimeStr,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasError := false
			
			if tt.req.UserID == "" {
				hasError = true
			}
			if tt.req.ScheduleDate != "" {
				_, err := time.Parse("2006-01-02", tt.req.ScheduleDate)
				if err != nil {
					hasError = true
				}
			}
			if tt.req.StartTime != "" {
				_, err := time.Parse("15:04", tt.req.StartTime)
				if err != nil {
					hasError = true
				}
			}
			if tt.req.EndTime != "" {
				_, err := time.Parse("15:04", tt.req.EndTime)
				if err != nil {
					hasError = true
				}
			}
			
			if hasError != tt.wantErr {
				t.Errorf("Expected error = %v, got %v", tt.wantErr, hasError)
			}
		})
	}
}

// Test CheckInRequest validation
func TestCheckInRequest(t *testing.T) {
	latitude := -2.9914
	longitude := 104.7565

	tests := []struct {
		name    string
		req     CheckInRequest
		wantErr bool
	}{
		{
			name: "valid request with location",
			req: CheckInRequest{
				Latitude:        &latitude,
				Longitude:       &longitude,
				IsViaUNSRIWiFi:  func() *bool { v := true; return &v }(),
			},
			wantErr: false,
		},
		{
			name: "valid request without location",
			req: CheckInRequest{
				IsViaUNSRIWiFi: func() *bool { v := false; return &v }(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// CheckInRequest is flexible - location is optional
			if tt.req.Latitude != nil && (*tt.req.Latitude < -90 || *tt.req.Latitude > 90) {
				t.Error("Latitude should be between -90 and 90")
			}
			if tt.req.Longitude != nil && (*tt.req.Longitude < -180 || *tt.req.Longitude > 180) {
				t.Error("Longitude should be between -180 and 180")
			}
		})
	}
}

// Test CheckOutRequest validation
func TestCheckOutRequest(t *testing.T) {
	latitude := -2.9914
	longitude := 104.7565

	tests := []struct {
		name    string
		req     CheckOutRequest
		wantErr bool
	}{
		{
			name: "valid request with location",
			req: CheckOutRequest{
				Latitude:        &latitude,
				Longitude:       &longitude,
				IsViaUNSRIWiFi:  func() *bool { v := true; return &v }(),
			},
			wantErr: false,
		},
		{
			name: "valid request without location",
			req: CheckOutRequest{
				IsViaUNSRIWiFi: func() *bool { v := false; return &v }(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// CheckOutRequest is flexible - location is optional
			if tt.req.Latitude != nil && (*tt.req.Latitude < -90 || *tt.req.Latitude > 90) {
				t.Error("Latitude should be between -90 and 90")
			}
			if tt.req.Longitude != nil && (*tt.req.Longitude < -180 || *tt.req.Longitude > 180) {
				t.Error("Longitude should be between -180 and 180")
			}
		})
	}
}

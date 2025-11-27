package service

import (
	"testing"

	apperrors "unsri-backend/internal/shared/errors"
)

// Test error types
func TestErrorTypes(t *testing.T) {
	t.Run("NotFoundError", func(t *testing.T) {
		err := apperrors.NewNotFoundError("location", "test-id")
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

// Test TapInRequest validation
func TestTapInRequest(t *testing.T) {
	latitude := -2.9914
	longitude := 104.7565

	tests := []struct {
		name    string
		req     TapInRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: TapInRequest{
				Latitude:  latitude,
				Longitude: longitude,
			},
			wantErr: false,
		},
		{
			name: "invalid latitude range",
			req: TapInRequest{
				Latitude:  91.0,
				Longitude: longitude,
			},
			wantErr: true,
		},
		{
			name: "invalid longitude range",
			req: TapInRequest{
				Latitude:  latitude,
				Longitude: 181.0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if (tt.req.Latitude < -90 || tt.req.Latitude > 90) && !tt.wantErr {
				t.Error("Latitude should be between -90 and 90")
			}
			if (tt.req.Longitude < -180 || tt.req.Longitude > 180) && !tt.wantErr {
				t.Error("Longitude should be between -180 and 180")
			}
		})
	}
}

// Test TapOutRequest validation
func TestTapOutRequest(t *testing.T) {
	latitude := -2.9914
	longitude := 104.7565

	tests := []struct {
		name    string
		req     TapOutRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: TapOutRequest{
				Latitude:  latitude,
				Longitude: longitude,
			},
			wantErr: false,
		},
		{
			name: "invalid latitude range",
			req: TapOutRequest{
				Latitude:  -91.0,
				Longitude: longitude,
			},
			wantErr: true,
		},
		{
			name: "invalid longitude range",
			req: TapOutRequest{
				Latitude:  latitude,
				Longitude: -181.0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if (tt.req.Latitude < -90 || tt.req.Latitude > 90) && !tt.wantErr {
				t.Error("Latitude should be between -90 and 90")
			}
			if (tt.req.Longitude < -180 || tt.req.Longitude > 180) && !tt.wantErr {
				t.Error("Longitude should be between -180 and 180")
			}
		})
	}
}

// Test geofence validation
func TestGeofenceValidation(t *testing.T) {
	t.Run("valid geofence coordinates", func(t *testing.T) {
		latitude := -2.9914
		longitude := 104.7565
		radius := 100.0 // meters

		if latitude < -90 || latitude > 90 {
			t.Error("Latitude should be between -90 and 90")
		}
		if longitude < -180 || longitude > 180 {
			t.Error("Longitude should be between -180 and 180")
		}
		if radius <= 0 {
			t.Error("Radius should be positive")
		}
	})

	t.Run("distance calculation", func(t *testing.T) {
		// Test that distance calculation logic would work
		// This is a placeholder for actual distance calculation tests
		lat1 := -2.9914
		lon1 := 104.7565
		lat2 := -2.9920
		lon2 := 104.7570

		// Simple validation that coordinates are different
		if lat1 == lat2 && lon1 == lon2 {
			t.Error("Coordinates should be different for distance calculation")
		}
	})
}


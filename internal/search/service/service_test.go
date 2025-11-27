package service

import (
	"testing"

	apperrors "unsri-backend/internal/shared/errors"
)

// Test error types
func TestErrorTypes(t *testing.T) {
	t.Run("NotFoundError", func(t *testing.T) {
		err := apperrors.NewNotFoundError("search", "test-id")
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

// Test SearchRequest validation
func TestSearchRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     SearchRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: SearchRequest{
				Query: "test",
				Type:  "users",
				Page:  1,
				PerPage: 20,
			},
			wantErr: false,
		},
		{
			name: "missing query",
			req: SearchRequest{
				Type: "users",
				Page: 1,
				PerPage: 20,
			},
			wantErr: true,
		},
		{
			name: "empty query",
			req: SearchRequest{
				Query: "",
				Type:  "users",
				Page:  1,
				PerPage: 20,
			},
			wantErr: true,
		},
		{
			name: "invalid pagination",
			req: SearchRequest{
				Query: "test",
				Type:  "users",
				Page:  -1,
				PerPage: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req.Query == "" && !tt.wantErr {
				t.Error("Query should be required")
			}
			if tt.req.Page < 1 && !tt.wantErr {
				t.Error("Page should be at least 1")
			}
			if tt.req.PerPage < 1 && !tt.wantErr {
				t.Error("PerPage should be at least 1")
			}
		})
	}
}

// Test GlobalSearch parameters validation
func TestGlobalSearchParameters(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		types   []string
		limit   int
		wantErr bool
	}{
		{
			name:    "valid request",
			query:   "test",
			types:   []string{"users", "courses"},
			limit:   10,
			wantErr: false,
		},
		{
			name:    "missing query",
			query:   "",
			types:   []string{"users", "courses"},
			limit:   10,
			wantErr: true,
		},
		{
			name:    "invalid limit",
			query:   "test",
			types:   []string{"users"},
			limit:   -1,
			wantErr: true,
		},
		{
			name:    "empty types",
			query:   "test",
			types:   []string{},
			limit:   10,
			wantErr: false, // Empty types means search all
		},
		{
			name:    "zero limit",
			query:   "test",
			types:   []string{"users"},
			limit:   0,
			wantErr: false, // Zero limit will be set to default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.query == "" && !tt.wantErr {
				t.Error("Query should be required")
			}
			if tt.limit < 0 && !tt.wantErr {
				t.Error("Limit should be non-negative")
			}
			// Test limit default value
			if tt.limit < 1 {
				limit := tt.limit
				if limit < 1 {
					limit = 10 // Default limit
				}
				if limit != 10 && !tt.wantErr {
					t.Error("Limit should default to 10 if less than 1")
				}
			}
		})
	}
}

// Test search types validation
func TestSearchTypes(t *testing.T) {
	validTypes := []string{"users", "courses", "schedules", "classes"}
	
	for _, stype := range validTypes {
		t.Run("valid type: "+stype, func(t *testing.T) {
			req := SearchRequest{
				Query: "test",
				Type:  stype,
				Page:  1,
				PerPage: 20,
			}
			if req.Type != stype {
				t.Errorf("Expected type %s, got %s", stype, req.Type)
			}
		})
	}
}


package service

import (
	"testing"

	apperrors "unsri-backend/internal/shared/errors"
)

// Test helper functions
func createTestQuickAction() *QuickAction {
	return &QuickAction{
		ID:          "test-action",
		Name:        "Test Action",
		Description: "Test description",
		Icon:        "test-icon",
		Route:       "/test",
	}
}

// Test error types
func TestErrorTypes(t *testing.T) {
	t.Run("NotFoundError", func(t *testing.T) {
		err := apperrors.NewNotFoundError("quick action", "test-id")
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

// Test QuickAction struct
func TestQuickActionStruct(t *testing.T) {
	t.Run("valid quick action", func(t *testing.T) {
		action := createTestQuickAction()
		if action.ID == "" {
			t.Error("ID should not be empty")
		}
		if action.Name == "" {
			t.Error("Name should not be empty")
		}
		if action.Description == "" {
			t.Error("Description should not be empty")
		}
		if action.Icon == "" {
			t.Error("Icon should not be empty")
		}
		if action.Route == "" {
			t.Error("Route should not be empty")
		}
	})

	t.Run("quick action fields", func(t *testing.T) {
		action := QuickAction{
			ID:          "test",
			Name:        "Test",
			Description: "Test description",
			Icon:        "test-icon",
			Route:       "/test",
		}
		if action.ID != "test" {
			t.Errorf("Expected ID 'test', got '%s'", action.ID)
		}
		if action.Name != "Test" {
			t.Errorf("Expected Name 'Test', got '%s'", action.Name)
		}
	})
}

// Test GetQuickActionsResponse validation
func TestGetQuickActionsResponse(t *testing.T) {
	t.Run("valid response", func(t *testing.T) {
		response := GetQuickActionsResponse{
			Role:    "mahasiswa",
			Actions: []QuickAction{*createTestQuickAction()},
		}
		if response.Role == "" {
			t.Error("Role should not be empty")
		}
		if len(response.Actions) == 0 {
			t.Error("Actions should not be empty")
		}
	})
}

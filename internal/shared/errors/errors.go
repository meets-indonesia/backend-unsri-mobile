package errors

import "fmt"

// AppError represents an application error
type AppError struct {
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s - %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// Predefined error codes
var (
	ErrCodeNotFound         = "NOT_FOUND"
	ErrCodeUnauthorized     = "UNAUTHORIZED"
	ErrCodeForbidden        = "FORBIDDEN"
	ErrCodeBadRequest       = "BAD_REQUEST"
	ErrCodeInternalError    = "INTERNAL_ERROR"
	ErrCodeValidationFailed = "VALIDATION_FAILED"
	ErrCodeConflict         = "CONFLICT"
)

// Error constructors
func NewNotFoundError(resource string, id string) *AppError {
	return &AppError{
		Code:    ErrCodeNotFound,
		Message: fmt.Sprintf("%s with id %s not found", resource, id),
	}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Code:    ErrCodeUnauthorized,
		Message: message,
	}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{
		Code:    ErrCodeForbidden,
		Message: message,
	}
}

func NewBadRequestError(message string) *AppError {
	return &AppError{
		Code:    ErrCodeBadRequest,
		Message: message,
	}
}

func NewInternalError(message string, err error) *AppError {
	return &AppError{
		Code:    ErrCodeInternalError,
		Message: message,
		Err:     err,
	}
}

func NewValidationError(message string) *AppError {
	return &AppError{
		Code:    ErrCodeValidationFailed,
		Message: message,
	}
}

func NewConflictError(message string) *AppError {
	return &AppError{
		Code:    ErrCodeConflict,
		Message: message,
	}
}


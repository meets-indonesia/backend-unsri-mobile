package utils

import (
	"fmt"
	"net/http"
	"strings"

	"unsri-backend/internal/shared/errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// ErrorInfo represents error information
type ErrorInfo struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Meta represents pagination or metadata
type Meta struct {
	Page       int   `json:"page,omitempty"`
	PerPage    int   `json:"per_page,omitempty"`
	Total      int64 `json:"total,omitempty"`
	TotalPages int   `json:"total_pages,omitempty"`
}

// SuccessResponse sends a success response
func SuccessResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, statusCode int, err error) {
	var appErr *errors.AppError
	var ok bool

	// Handle validation errors from gin
	if validationErr, ok := err.(validator.ValidationErrors); ok {
		details := make(map[string]interface{})
		var messages []string

		for _, fieldErr := range validationErr {
			field := strings.ToLower(fieldErr.Field())
			tag := fieldErr.Tag()

			var message string
			switch tag {
			case "required":
				message = fmt.Sprintf("%s is required", field)
			case "email":
				message = fmt.Sprintf("%s must be a valid email", field)
			case "min":
				message = fmt.Sprintf("%s must be at least %s characters", field, fieldErr.Param())
			case "max":
				message = fmt.Sprintf("%s must be at most %s characters", field, fieldErr.Param())
			case "oneof":
				message = fmt.Sprintf("%s must be one of: %s", field, fieldErr.Param())
			default:
				message = fmt.Sprintf("%s is invalid", field)
			}

			details[field] = message
			messages = append(messages, message)
		}

		appErr = &errors.AppError{
			Code:    errors.ErrCodeValidationFailed,
			Message: "Validation failed: " + strings.Join(messages, "; "),
		}

		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error: &ErrorInfo{
				Code:    appErr.Code,
				Message: appErr.Message,
				Details: details,
			},
		})
		return
	}

	// Handle gin binding errors
	if ginErr, ok := err.(*gin.Error); ok {
		if validationErr, ok := ginErr.Err.(validator.ValidationErrors); ok {
			ErrorResponse(c, statusCode, validationErr)
			return
		}
	}

	// Handle standard errors
	if appErr, ok = err.(*errors.AppError); !ok {
		// Convert to AppError with original error message
		errorMessage := "Internal server error"
		if err != nil && err.Error() != "" {
			errorMessage = err.Error()
		}
		appErr = errors.NewInternalError(errorMessage, err)
		statusCode = http.StatusInternalServerError
	}

	// Map error codes to HTTP status codes
	// If statusCode is 0, determine it from error code
	if statusCode == 0 {
		switch appErr.Code {
		case errors.ErrCodeNotFound:
			statusCode = http.StatusNotFound
		case errors.ErrCodeUnauthorized:
			statusCode = http.StatusUnauthorized
		case errors.ErrCodeForbidden:
			statusCode = http.StatusForbidden
		case errors.ErrCodeBadRequest, errors.ErrCodeValidationFailed:
			statusCode = http.StatusBadRequest
		case errors.ErrCodeConflict:
			statusCode = http.StatusConflict
		case errors.ErrCodeInternalError:
			statusCode = http.StatusInternalServerError
		default:
			statusCode = http.StatusInternalServerError
		}
	} else {
		// If statusCode is provided, still map error codes to ensure consistency
		if appErr.Code == errors.ErrCodeNotFound {
			statusCode = http.StatusNotFound
		} else if appErr.Code == errors.ErrCodeUnauthorized {
			statusCode = http.StatusUnauthorized
		} else if appErr.Code == errors.ErrCodeForbidden {
			statusCode = http.StatusForbidden
		} else if appErr.Code == errors.ErrCodeBadRequest || appErr.Code == errors.ErrCodeValidationFailed {
			statusCode = http.StatusBadRequest
		} else if appErr.Code == errors.ErrCodeConflict {
			statusCode = http.StatusConflict
		} else if appErr.Code == errors.ErrCodeInternalError {
			statusCode = http.StatusInternalServerError
		}
	}

	// Build error response with details
	errorInfo := &ErrorInfo{
		Code:    appErr.Code,
		Message: appErr.Message,
	}

	// Include error details when underlying error exists
	if appErr.Err != nil {
		details := make(map[string]interface{})
		details["error"] = appErr.Err.Error()

		// For internal errors, try to extract more specific error information
		if appErr.Code == errors.ErrCodeInternalError {
			if unwrappedErr := appErr.Unwrap(); unwrappedErr != nil && unwrappedErr != appErr.Err {
				details["underlying_error"] = unwrappedErr.Error()
			}
		}

		errorInfo.Details = details
	}

	c.JSON(statusCode, Response{
		Success: false,
		Error:   errorInfo,
	})
}

// ValidationErrorResponse sends a validation error response with field details
func ValidationErrorResponse(c *gin.Context, err error) {
	ErrorResponse(c, http.StatusBadRequest, err)
}

// BadRequestResponse sends a bad request error response
func BadRequestResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusBadRequest, errors.NewBadRequestError(message))
}

// NotFoundResponse sends a not found error response
func NotFoundResponse(c *gin.Context, resource string, id string) {
	ErrorResponse(c, http.StatusNotFound, errors.NewNotFoundError(resource, id))
}

// UnauthorizedResponse sends an unauthorized error response
func UnauthorizedResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, errors.NewUnauthorizedError(message))
}

// ForbiddenResponse sends a forbidden error response
func ForbiddenResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusForbidden, errors.NewForbiddenError(message))
}

// ConflictResponse sends a conflict error response
func ConflictResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusConflict, errors.NewConflictError(message))
}

// InternalErrorResponse sends an internal server error response
func InternalErrorResponse(c *gin.Context, message string, err error) {
	ErrorResponse(c, http.StatusInternalServerError, errors.NewInternalError(message, err))
}

// PaginatedResponse sends a paginated response
func PaginatedResponse(c *gin.Context, data interface{}, page, perPage int, total int64) {
	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
		Meta: &Meta{
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

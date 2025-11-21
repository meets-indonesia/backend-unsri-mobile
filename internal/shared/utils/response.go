package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"unsri-backend/internal/shared/errors"
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
	Code    string `json:"code"`
	Message string `json:"message"`
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

	if appErr, ok = err.(*errors.AppError); !ok {
		appErr = errors.NewInternalError("Internal server error", err)
		statusCode = http.StatusInternalServerError
	}

	// Map error codes to HTTP status codes
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
	}

	c.JSON(statusCode, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    appErr.Code,
			Message: appErr.Message,
		},
	})
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


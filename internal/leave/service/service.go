package service

import (
	"context"
	"time"

	"unsri-backend/internal/leave/repository"
	apperrors "unsri-backend/internal/shared/errors"
	"unsri-backend/internal/shared/models"
)

// LeaveService handles leave business logic
type LeaveService struct {
	repo *repository.LeaveRepository
}

// NewLeaveService creates a new leave service
func NewLeaveService(repo *repository.LeaveRepository) *LeaveService {
	return &LeaveService{repo: repo}
}

// ========== Leave Request Service Methods ==========

// CreateLeaveRequestRequest represents create leave request
type CreateLeaveRequestRequest struct {
	LeaveType     string  `json:"leave_type" binding:"required,oneof=ANNUAL_LEAVE SICK_LEAVE PERSONAL_LEAVE EMERGENCY_LEAVE UNPAID_LEAVE OTHER"`
	StartDate     string  `json:"start_date" binding:"required"`
	EndDate       string  `json:"end_date" binding:"required"`
	Reason        string  `json:"reason" binding:"required"`
	AttachmentURL *string `json:"attachment_url,omitempty"`
}

// CreateLeaveRequest creates a new leave request
func (s *LeaveService) CreateLeaveRequest(ctx context.Context, userID string, req CreateLeaveRequestRequest) (*models.LeaveRequest, error) {
	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, apperrors.NewValidationError("invalid start_date format, use YYYY-MM-DD")
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, apperrors.NewValidationError("invalid end_date format, use YYYY-MM-DD")
	}

	// Validate dates
	if endDate.Before(startDate) {
		return nil, apperrors.NewValidationError("end_date must be after start_date")
	}

	// Calculate total days (inclusive of start and end date)
	totalDays := endDate.Sub(startDate).Hours()/24 + 1

	// Check leave quota for quota-based leave types
	leaveType := models.LeaveType(req.LeaveType)
	if leaveType == models.LeaveTypeAnnual || leaveType == models.LeaveTypeSick {
		currentYear := time.Now().Year()
		quota, err := s.repo.GetLeaveQuotaByUserAndTypeAndYear(ctx, userID, leaveType, currentYear)
		if err != nil {
			return nil, apperrors.NewInternalError("failed to check leave quota", err)
		}

		if quota.RemainingQuota < totalDays {
			return nil, apperrors.NewValidationError("insufficient leave quota")
		}
	}

	leaveRequest := &models.LeaveRequest{
		UserID:        userID,
		LeaveType:     leaveType,
		StartDate:     startDate,
		EndDate:       endDate,
		TotalDays:     totalDays,
		Reason:        req.Reason,
		Status:        models.LeaveStatusPending,
		AttachmentURL: req.AttachmentURL,
	}

	if err := s.repo.CreateLeaveRequest(ctx, leaveRequest); err != nil {
		return nil, apperrors.NewInternalError("failed to create leave request", err)
	}

	return leaveRequest, nil
}

// GetLeaveRequestByID gets a leave request by ID
func (s *LeaveService) GetLeaveRequestByID(ctx context.Context, id string) (*models.LeaveRequest, error) {
	leaveRequest, err := s.repo.GetLeaveRequestByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewNotFoundError("leave request", id)
	}
	return leaveRequest, nil
}

// GetLeaveRequestsRequest represents get leave requests request
type GetLeaveRequestsRequest struct {
	UserID    string `form:"user_id"`
	LeaveType string `form:"leave_type"`
	Status    string `form:"status"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Page      int    `form:"page,default=1"`
	PerPage   int    `form:"per_page,default=20"`
}

// GetLeaveRequests gets all leave requests
func (s *LeaveService) GetLeaveRequests(ctx context.Context, req GetLeaveRequestsRequest) ([]models.LeaveRequest, int64, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}

	var userIDPtr, leaveTypePtr, statusPtr *string
	if req.UserID != "" {
		userIDPtr = &req.UserID
	}
	if req.LeaveType != "" {
		leaveTypePtr = &req.LeaveType
	}
	if req.Status != "" {
		statusPtr = &req.Status
	}

	var startDatePtr, endDatePtr *time.Time
	if req.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err == nil {
			startDatePtr = &startDate
		}
	}
	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err == nil {
			endDatePtr = &endDate
		}
	}

	return s.repo.GetAllLeaveRequests(ctx, userIDPtr, leaveTypePtr, statusPtr, startDatePtr, endDatePtr, perPage, (page-1)*perPage)
}

// GetLeaveRequestsByUser gets leave requests by user ID
func (s *LeaveService) GetLeaveRequestsByUser(ctx context.Context, userID string, status *string) ([]models.LeaveRequest, error) {
	return s.repo.GetLeaveRequestsByUserID(ctx, userID, status)
}

// ApproveLeaveRequestRequest represents approve leave request
type ApproveLeaveRequestRequest struct {
	Notes string `json:"notes,omitempty"`
}

// ApproveLeaveRequest approves a leave request
func (s *LeaveService) ApproveLeaveRequest(ctx context.Context, leaveID string, approverID string, req ApproveLeaveRequestRequest) (*models.LeaveRequest, error) {
	leaveRequest, err := s.repo.GetLeaveRequestByID(ctx, leaveID)
	if err != nil {
		return nil, apperrors.NewNotFoundError("leave request", leaveID)
	}

	if leaveRequest.Status != models.LeaveStatusPending {
		return nil, apperrors.NewValidationError("only pending leave requests can be approved")
	}

	now := time.Now()
	leaveRequest.Status = models.LeaveStatusApproved
	leaveRequest.ApprovedBy = &approverID
	leaveRequest.ApprovedAt = &now

	// Update leave quota if applicable
	if leaveRequest.LeaveType == models.LeaveTypeAnnual || leaveRequest.LeaveType == models.LeaveTypeSick {
		currentYear := time.Now().Year()
		quota, err := s.repo.GetLeaveQuotaByUserAndTypeAndYear(ctx, leaveRequest.UserID, leaveRequest.LeaveType, currentYear)
		if err != nil {
			// Log error but continue
			_ = err
		} else {
			quota.UsedQuota += leaveRequest.TotalDays
			quota.RemainingQuota = quota.TotalQuota - quota.UsedQuota
			if err := s.repo.UpdateLeaveQuota(ctx, quota); err != nil {
				// Log error but continue
				_ = err
			}
		}
	}

	if err := s.repo.UpdateLeaveRequest(ctx, leaveRequest); err != nil {
		return nil, apperrors.NewInternalError("failed to approve leave request", err)
	}

	return leaveRequest, nil
}

// RejectLeaveRequestRequest represents reject leave request
type RejectLeaveRequestRequest struct {
	RejectionReason string `json:"rejection_reason" binding:"required"`
}

// RejectLeaveRequest rejects a leave request
func (s *LeaveService) RejectLeaveRequest(ctx context.Context, leaveID string, approverID string, req RejectLeaveRequestRequest) (*models.LeaveRequest, error) {
	leaveRequest, err := s.repo.GetLeaveRequestByID(ctx, leaveID)
	if err != nil {
		return nil, apperrors.NewNotFoundError("leave request", leaveID)
	}

	if leaveRequest.Status != models.LeaveStatusPending {
		return nil, apperrors.NewValidationError("only pending leave requests can be rejected")
	}

	leaveRequest.Status = models.LeaveStatusRejected
	leaveRequest.ApprovedBy = &approverID
	rejectionReason := req.RejectionReason
	leaveRequest.RejectionReason = &rejectionReason

	if err := s.repo.UpdateLeaveRequest(ctx, leaveRequest); err != nil {
		return nil, apperrors.NewInternalError("failed to reject leave request", err)
	}

	return leaveRequest, nil
}

// CancelLeaveRequest cancels a leave request (by user)
func (s *LeaveService) CancelLeaveRequest(ctx context.Context, leaveID string, userID string) (*models.LeaveRequest, error) {
	leaveRequest, err := s.repo.GetLeaveRequestByID(ctx, leaveID)
	if err != nil {
		return nil, apperrors.NewNotFoundError("leave request", leaveID)
	}

	// Only owner can cancel
	if leaveRequest.UserID != userID {
		return nil, apperrors.NewForbiddenError("only leave request owner can cancel")
	}

	if leaveRequest.Status != models.LeaveStatusPending {
		return nil, apperrors.NewValidationError("only pending leave requests can be cancelled")
	}

	// If already approved, revert quota before cancelling
	if leaveRequest.Status == models.LeaveStatusApproved {
		if leaveRequest.LeaveType == models.LeaveTypeAnnual || leaveRequest.LeaveType == models.LeaveTypeSick {
			currentYear := time.Now().Year()
			quota, err := s.repo.GetLeaveQuotaByUserAndTypeAndYear(ctx, leaveRequest.UserID, leaveRequest.LeaveType, currentYear)
			if err != nil {
				// Log error but continue
				_ = err
			} else {
				quota.UsedQuota -= leaveRequest.TotalDays
				if quota.UsedQuota < 0 {
					quota.UsedQuota = 0
				}
				quota.RemainingQuota = quota.TotalQuota - quota.UsedQuota
				if err := s.repo.UpdateLeaveQuota(ctx, quota); err != nil {
					// Log error but continue
					_ = err
				}
			}
		}
	}

	leaveRequest.Status = models.LeaveStatusCancelled

	if err := s.repo.UpdateLeaveRequest(ctx, leaveRequest); err != nil {
		return nil, apperrors.NewInternalError("failed to cancel leave request", err)
	}

	return leaveRequest, nil
}

// DeleteLeaveRequest deletes a leave request
func (s *LeaveService) DeleteLeaveRequest(ctx context.Context, id string) error {
	_, err := s.repo.GetLeaveRequestByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFoundError("leave request", id)
	}
	return s.repo.DeleteLeaveRequest(ctx, id)
}

// ========== Leave Quota Service Methods ==========

// CreateLeaveQuotaRequest represents create leave quota request
type CreateLeaveQuotaRequest struct {
	UserID     string  `json:"user_id" binding:"required"`
	LeaveType  string  `json:"leave_type" binding:"required,oneof=ANNUAL_LEAVE SICK_LEAVE PERSONAL_LEAVE EMERGENCY_LEAVE UNPAID_LEAVE OTHER"`
	Year       int     `json:"year" binding:"required"`
	TotalQuota float64 `json:"total_quota" binding:"required"`
}

// CreateLeaveQuota creates a new leave quota
func (s *LeaveService) CreateLeaveQuota(ctx context.Context, req CreateLeaveQuotaRequest) (*models.LeaveQuota, error) {
	// Check if quota already exists
	_, err := s.repo.GetLeaveQuotaByUserAndTypeAndYear(ctx, req.UserID, models.LeaveType(req.LeaveType), req.Year)
	if err == nil {
		return nil, apperrors.NewConflictError("leave quota already exists for this user, type, and year")
	}

	quota := &models.LeaveQuota{
		UserID:     req.UserID,
		LeaveType:  models.LeaveType(req.LeaveType),
		Year:       req.Year,
		TotalQuota: req.TotalQuota,
		UsedQuota:  0,
	}

	if err := s.repo.CreateLeaveQuota(ctx, quota); err != nil {
		return nil, apperrors.NewInternalError("failed to create leave quota", err)
	}

	return quota, nil
}

// GetLeaveQuotaByID gets a leave quota by ID
func (s *LeaveService) GetLeaveQuotaByID(ctx context.Context, id string) (*models.LeaveQuota, error) {
	quota, err := s.repo.GetLeaveQuotaByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewNotFoundError("leave quota", id)
	}
	return quota, nil
}

// GetLeaveQuotasRequest represents get leave quotas request
type GetLeaveQuotasRequest struct {
	UserID    string `form:"user_id"`
	LeaveType string `form:"leave_type"`
	Year      int    `form:"year"`
	Page      int    `form:"page,default=1"`
	PerPage   int    `form:"per_page,default=20"`
}

// GetLeaveQuotas gets all leave quotas
func (s *LeaveService) GetLeaveQuotas(ctx context.Context, req GetLeaveQuotasRequest) ([]models.LeaveQuota, int64, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 20
	}

	var userIDPtr, leaveTypePtr *string
	if req.UserID != "" {
		userIDPtr = &req.UserID
	}
	if req.LeaveType != "" {
		leaveTypePtr = &req.LeaveType
	}

	var yearPtr *int
	if req.Year > 0 {
		yearPtr = &req.Year
	}

	return s.repo.GetAllLeaveQuotas(ctx, userIDPtr, leaveTypePtr, yearPtr, perPage, (page-1)*perPage)
}

// GetLeaveQuotasByUser gets leave quotas by user ID
func (s *LeaveService) GetLeaveQuotasByUser(ctx context.Context, userID string, year *int) ([]models.LeaveQuota, error) {
	return s.repo.GetLeaveQuotasByUserID(ctx, userID, year)
}

// UpdateLeaveQuotaRequest represents update leave quota request
type UpdateLeaveQuotaRequest struct {
	TotalQuota *float64 `json:"total_quota,omitempty"`
	UsedQuota  *float64 `json:"used_quota,omitempty"`
}

// UpdateLeaveQuota updates a leave quota
func (s *LeaveService) UpdateLeaveQuota(ctx context.Context, id string, req UpdateLeaveQuotaRequest) (*models.LeaveQuota, error) {
	quota, err := s.repo.GetLeaveQuotaByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewNotFoundError("leave quota", id)
	}

	if req.TotalQuota != nil {
		quota.TotalQuota = *req.TotalQuota
	}
	if req.UsedQuota != nil {
		quota.UsedQuota = *req.UsedQuota
	}

	// Recalculate remaining quota
	quota.RemainingQuota = quota.TotalQuota - quota.UsedQuota

	if err := s.repo.UpdateLeaveQuota(ctx, quota); err != nil {
		return nil, apperrors.NewInternalError("failed to update leave quota", err)
	}

	return quota, nil
}

// DeleteLeaveQuota deletes a leave quota
func (s *LeaveService) DeleteLeaveQuota(ctx context.Context, id string) error {
	_, err := s.repo.GetLeaveQuotaByID(ctx, id)
	if err != nil {
		return apperrors.NewNotFoundError("leave quota", id)
	}
	return s.repo.DeleteLeaveQuota(ctx, id)
}

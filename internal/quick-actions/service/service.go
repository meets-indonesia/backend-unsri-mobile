package service

import (
	"context"

	apperrors "unsri-backend/internal/shared/errors"
	"unsri-backend/internal/shared/models"
	"unsri-backend/internal/quick-actions/repository"
)

// QuickActionsService handles quick actions business logic
type QuickActionsService struct {
	repo *repository.QuickActionsRepository
}

// NewQuickActionsService creates a new quick actions service
func NewQuickActionsService(repo *repository.QuickActionsRepository) *QuickActionsService {
	return &QuickActionsService{repo: repo}
}

// GetQuickActionsResponse represents quick actions response
type GetQuickActionsResponse struct {
	Role        string                   `json:"role"`
	Actions     []QuickAction            `json:"actions"`
}

// QuickAction represents a quick action
type QuickAction struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Route       string `json:"route"`
}

// GetQuickActions gets available quick actions by role
func (s *QuickActionsService) GetQuickActions(ctx context.Context, role string) (*GetQuickActionsResponse, error) {
	actions := []QuickAction{}

	switch role {
	case "mahasiswa":
		actions = []QuickAction{
			{ID: "transcript", Name: "Transkrip", Description: "Lihat transkrip nilai", Icon: "transcript", Route: "/transcript"},
			{ID: "krs", Name: "KRS", Description: "Kartu Rencana Studi", Icon: "krs", Route: "/krs"},
			{ID: "bimbingan", Name: "Bimbingan", Description: "Riwayat bimbingan", Icon: "bimbingan", Route: "/bimbingan"},
		}
	case "dosen":
		actions = []QuickAction{
			{ID: "bimbingan", Name: "Bimbingan", Description: "Daftar bimbingan", Icon: "bimbingan", Route: "/bimbingan"},
			{ID: "payroll", Name: "Payroll", Description: "Informasi gaji", Icon: "payroll", Route: "/payroll"},
			{ID: "claims", Name: "Klaim", Description: "Klaim pengeluaran", Icon: "claims", Route: "/claims"},
			{ID: "leaves", Name: "Cuti", Description: "Pengajuan cuti", Icon: "leaves", Route: "/leaves"},
		}
	case "staff":
		actions = []QuickAction{
			{ID: "payroll", Name: "Payroll", Description: "Informasi gaji", Icon: "payroll", Route: "/payroll"},
			{ID: "claims", Name: "Klaim", Description: "Klaim pengeluaran", Icon: "claims", Route: "/claims"},
			{ID: "leaves", Name: "Cuti", Description: "Pengajuan cuti", Icon: "leaves", Route: "/leaves"},
		}
	}

	return &GetQuickActionsResponse{
		Role:    role,
		Actions: actions,
	}, nil
}

// GetTranscript gets transcript for a student
func (s *QuickActionsService) GetTranscript(ctx context.Context, studentID string) (*models.Transcript, error) {
	return s.repo.GetTranscriptByStudentID(ctx, studentID)
}

// GetKRS gets KRS for a student
func (s *QuickActionsService) GetKRS(ctx context.Context, studentID string, semester *string) (*models.KRS, error) {
	return s.repo.GetKRSByStudentID(ctx, studentID, semester)
}

// GetBimbingans gets bimbingan records
func (s *QuickActionsService) GetBimbingans(ctx context.Context, userID string, role string, limit int) ([]models.Bimbingan, error) {
	if role == "mahasiswa" {
		return s.repo.GetBimbingansByStudentID(ctx, userID, limit)
	} else if role == "dosen" {
		return s.repo.GetBimbingansByDosenID(ctx, userID, limit)
	}
	return nil, apperrors.NewBadRequestError("invalid role for bimbingan")
}


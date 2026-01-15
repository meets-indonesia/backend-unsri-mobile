package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LeaveType represents leave type
type LeaveType string

const (
	LeaveTypeAnnual    LeaveType = "ANNUAL_LEAVE"
	LeaveTypeSick      LeaveType = "SICK_LEAVE"
	LeaveTypePersonal  LeaveType = "PERSONAL_LEAVE"
	LeaveTypeEmergency LeaveType = "EMERGENCY_LEAVE"
	LeaveTypeUnpaid    LeaveType = "UNPAID_LEAVE"
	LeaveTypeOther     LeaveType = "OTHER"
)

// LeaveStatus represents leave request status
type LeaveStatus string

const (
	LeaveStatusPending   LeaveStatus = "PENDING"
	LeaveStatusApproved  LeaveStatus = "APPROVED"
	LeaveStatusRejected  LeaveStatus = "REJECTED"
	LeaveStatusCancelled LeaveStatus = "CANCELLED"
)

// LeaveRequest represents a leave request
type LeaveRequest struct {
	ID              string     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID          string     `gorm:"type:uuid;not null;index" json:"user_id"`
	LeaveType       LeaveType  `gorm:"type:varchar(20);not null" json:"leave_type"`
	StartDate       time.Time  `gorm:"type:date;not null" json:"start_date"`
	EndDate         time.Time  `gorm:"type:date;not null" json:"end_date"`
	TotalDays       float64    `gorm:"type:decimal(5,2);not null" json:"total_days"`
	Reason          string     `gorm:"type:text;not null" json:"reason"`
	Status          LeaveStatus `gorm:"type:varchar(20);not null;default:'PENDING'" json:"status"`
	ApprovedBy      *string    `gorm:"type:uuid;index" json:"approved_by,omitempty"`
	ApprovedAt      *time.Time `gorm:"type:timestamp" json:"approved_at,omitempty"`
	RejectionReason *string    `gorm:"type:text" json:"rejection_reason,omitempty"`
	AttachmentURL   *string   `gorm:"type:varchar(500)" json:"attachment_url,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User      User  `gorm:"foreignKey:UserID" json:"-"`
	Approver  *User `gorm:"foreignKey:ApprovedBy" json:"-"`
}

// TableName specifies the table name
func (LeaveRequest) TableName() string {
	return "leave_requests"
}

// BeforeCreate hook
func (l *LeaveRequest) BeforeCreate(tx *gorm.DB) error {
	if l.ID == "" {
		l.ID = uuid.New().String()
	}
	return nil
}

// LeaveQuota represents leave quota for a user
type LeaveQuota struct {
	ID              string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID          string    `gorm:"type:uuid;not null;index" json:"user_id"`
	LeaveType       LeaveType `gorm:"type:varchar(20);not null" json:"leave_type"`
	Year            int       `gorm:"type:integer;not null" json:"year"`
	TotalQuota      float64   `gorm:"type:decimal(5,2);not null" json:"total_quota"`
	UsedQuota       float64   `gorm:"type:decimal(5,2);default:0" json:"used_quota"`
	RemainingQuota  float64   `gorm:"type:decimal(5,2);default:0" json:"remaining_quota"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName specifies the table name
func (LeaveQuota) TableName() string {
	return "leave_quotas"
}

// BeforeCreate hook
func (l *LeaveQuota) BeforeCreate(tx *gorm.DB) error {
	if l.ID == "" {
		l.ID = uuid.New().String()
	}
	// Calculate remaining quota
	l.RemainingQuota = l.TotalQuota - l.UsedQuota
	return nil
}

// BeforeUpdate hook
func (l *LeaveQuota) BeforeUpdate(tx *gorm.DB) error {
	// Recalculate remaining quota
	l.RemainingQuota = l.TotalQuota - l.UsedQuota
	return nil
}

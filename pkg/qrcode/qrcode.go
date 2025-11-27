package qrcode

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/skip2/go-qrcode"
)

// QRData represents QR code data structure
type QRData struct {
	SessionID  string    `json:"session_id"`
	ScheduleID string    `json:"schedule_id"`
	ExpiresAt  time.Time `json:"expires_at"`
	Type       string    `json:"type"` // "kelas", "kampus", "gate"

	// Gate access specific fields (for UNSRI gate integration)
	UserID   string `json:"user_id,omitempty"`
	QRToken  string `json:"qr_token,omitempty"`
	UserName string `json:"user_name,omitempty"`
	UserRole string `json:"user_role,omitempty"`
	NIM      string `json:"nim,omitempty"`   // For mahasiswa
	NIP      string `json:"nip,omitempty"`   // For dosen/staff
	Prodi    string `json:"prodi,omitempty"` // Program Studi
}

// GenerateQRCode generates a QR code image from data
func GenerateQRCode(data QRData) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal QR data: %w", err)
	}

	png, err := qrcode.Encode(string(jsonData), qrcode.Medium, 256)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	return png, nil
}

// ParseQRData parses QR code data from string
func ParseQRData(data string) (*QRData, error) {
	var qrData QRData
	if err := json.Unmarshal([]byte(data), &qrData); err != nil {
		return nil, fmt.Errorf("failed to parse QR data: %w", err)
	}

	return &qrData, nil
}

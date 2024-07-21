// internal/model/transcribe.go
package model

type TranscriptionStatus string

const (
	StatusPending    TranscriptionStatus = "pending"
	StatusProcessing TranscriptionStatus = "processing"
	StatusCompleted  TranscriptionStatus = "completed"
	StatusFailed     TranscriptionStatus = "failed"
)

type TranscriptionTask struct {
	ID            string              `json:"id"`
	UserID        string              `json:"user_id"`
	Status        TranscriptionStatus `json:"status"`
	Transcription string              `json:"transcription,omitempty"`
	Error         string              `json:"error,omitempty"`
	CreatedAt     int64               `json:"created_at"`
	UpdatedAt     int64               `json:"updated_at"`
}

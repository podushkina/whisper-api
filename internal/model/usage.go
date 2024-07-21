// internal/model/usage.go
package model

import "time"

type Usage struct {
	UserID        string    `json:"user_id"`
	RequestCount  int       `json:"request_count"`
	TotalDuration int       `json:"total_duration"`
	TotalDataSize int64     `json:"total_data_size"`
	LastUsedAt    time.Time `json:"last_used_at"`
}

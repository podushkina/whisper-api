// internal/service/usage.go
package service

import (
	"whisper-api/internal/model"
	"whisper-api/internal/repository"
)

type UsageService struct {
	repo *repository.UsageRepository
}

func NewUsageService(repo *repository.UsageRepository) *UsageService {
	return &UsageService{repo: repo}
}

func (s *UsageService) GetUsage(userID string) (*model.Usage, error) {
	return s.repo.GetUsage(userID)
}

func (s *UsageService) IncrementUsage(userID string, duration int, dataSize int64) error {
	return s.repo.IncrementUsage(userID, duration, dataSize)
}

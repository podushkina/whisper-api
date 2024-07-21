// internal/repository/usage.go
package repository

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"

	"whisper-api/internal/model"
)

type UsageRepository struct {
	redis *redis.Client
}

func NewUsageRepository(redis *redis.Client) *UsageRepository {
	return &UsageRepository{redis: redis}
}

func (r *UsageRepository) GetUsage(userID string) (*model.Usage, error) {
	ctx := context.Background()
	usageJSON, err := r.redis.Get(ctx, "usage:"+userID).Bytes()
	if err != nil {
		if err == redis.Nil {
			return &model.Usage{UserID: userID}, nil
		} // TODO: попробовать error.Is
		return nil, err
	}

	var usage model.Usage
	err = json.Unmarshal(usageJSON, &usage)
	if err != nil {
		return nil, err
	}

	return &usage, nil
}

func (r *UsageRepository) IncrementUsage(userID string, duration int, dataSize int64) error {
	ctx := context.Background()
	usage, err := r.GetUsage(userID)
	if err != nil {
		return err
	}

	usage.RequestCount++
	usage.TotalDuration += duration
	usage.TotalDataSize += dataSize

	usageJSON, err := json.Marshal(usage)
	if err != nil {
		return err
	}

	return r.redis.Set(ctx, "usage:"+userID, usageJSON, 0).Err()
}

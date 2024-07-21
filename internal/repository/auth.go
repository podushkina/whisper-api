// internal/repository/auth.go
package repository

import (
	"context"
	"github.com/go-redis/redis/v8" // TODO: уже вышла сервия v9 переписать под нее
	"time"
)

type AuthRepository struct {
	redis *redis.Client
}

func NewAuthRepository(redis *redis.Client) *AuthRepository {
	return &AuthRepository{redis: redis}
}

func (r *AuthRepository) SaveToken(userID, token string) error {
	ctx := context.Background()
	return r.redis.Set(ctx, "token:"+userID, token, 24*time.Hour).Err()
}

func (r *AuthRepository) GetToken(userID string) (string, error) {
	ctx := context.Background()
	return r.redis.Get(ctx, "token:"+userID).Result()
}

// internal/service/auth.go
package service

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
	"whisper-api/internal/config"
	"whisper-api/internal/repository"
)

var ErrInvalidAPIKey = errors.New("invalid API key")

type AuthService struct {
	cfg  *config.Config
	repo *repository.AuthRepository
}

func NewAuthService(cfg *config.Config, repo *repository.AuthRepository) *AuthService {
	return &AuthService{cfg: cfg, repo: repo}
}

func (s *AuthService) GenerateToken(apiKey string) (string, error) {
	// TODO: в реальном приложении здесь должна быть проверка API ключа, добавить функцию ValidateToken

	userID := uuid.New().String()
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return "", err
	}

	err = s.repo.SaveToken(userID, tokenString)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"whisper-api/internal/config"
	"whisper-api/internal/handler"
	"whisper-api/internal/repository"
	"whisper-api/internal/server"
	"whisper-api/internal/service"
	"whisper-api/pkg/logger"
	"whisper-api/pkg/whisper"
)

// TODO: добавить тесты с использованием testify, go-playground/validator/v10, httpexpect и gofakeit
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	appLogger := logger.NewLogger()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	authRepo := repository.NewAuthRepository(redisClient)
	transcribeRepo := repository.NewTranscribeRepository(redisClient)
	usageRepo := repository.NewUsageRepository(redisClient)

	whisperClient := whisper.NewClient(cfg.Whisper.Path)

	authService := service.NewAuthService(cfg, authRepo)
	transcribeService := service.NewTranscribeService(transcribeRepo, whisperClient)
	usageService := service.NewUsageService(usageRepo)

	authHandler := handler.NewAuthHandler(authService)
	transcribeHandler := handler.NewTranscribeHandler(transcribeService)
	usageHandler := handler.NewUsageHandler(usageService)
	healthHandler := handler.NewHealthHandler()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Post("/auth/token", authHandler.GetToken)

	r.Group(func(r chi.Router) {
		r.Use(server.Auth(cfg))
		r.Post("/transcribe", transcribeHandler.CreateTranscription)
		r.Get("/transcribe", transcribeHandler.ListTranscriptions)
		r.Get("/transcribe/{id}", transcribeHandler.GetTranscription)
		r.Delete("/transcribe/{id}", transcribeHandler.DeleteTranscription)
		r.Get("/usage", usageHandler.GetUsage)
	})

	r.Get("/health", healthHandler.Check)

	srv := server.NewServer(cfg, appLogger, r)

	go func() {
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			appLogger.Error("Failed to run server", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Error("Server forced to shutdown", "error", err)
	}

	appLogger.Info("Server exiting")
}

package service

import (
	"io"

	"github.com/google/uuid"
	"whisper-api/internal/model"
	"whisper-api/internal/repository"
	"whisper-api/pkg/whisper"
)

type TranscribeService struct {
	repo          *repository.TranscribeRepository
	whisperClient *whisper.Client
}

func NewTranscribeService(repo *repository.TranscribeRepository, whisperClient *whisper.Client) *TranscribeService {
	return &TranscribeService{repo: repo, whisperClient: whisperClient}
}

func (s *TranscribeService) CreateTranscriptionTask(userID string, audioFile io.Reader, filename string, options whisper.TranscriptionOptions) (*model.TranscriptionTask, error) {
	taskID := uuid.New().String()
	task := &model.TranscriptionTask{
		ID:     taskID,
		UserID: userID,
		Status: model.StatusPending,
	}

	err := s.repo.SaveTask(task)
	if err != nil {
		return nil, err
	}

	go s.processTranscription(task, audioFile, options)

	return task, nil
}

func (s *TranscribeService) processTranscription(task *model.TranscriptionTask, audioFile io.Reader, options whisper.TranscriptionOptions) {
	task.Status = model.StatusProcessing
	s.repo.UpdateTask(task) // TODO: написать хендерл для ошибки

	transcription, err := s.whisperClient.Transcribe(audioFile, options)
	if err != nil {
		task.Status = model.StatusFailed
		task.Error = err.Error()
	} else {
		task.Status = model.StatusCompleted
		task.Transcription = transcription
	}

	s.repo.UpdateTask(task) // TODO: написать хендерл для ошибки
}

func (s *TranscribeService) GetTranscriptionTask(userID, taskID string) (*model.TranscriptionTask, error) {
	return s.repo.GetTask(userID, taskID)
}

func (s *TranscribeService) ListTranscriptionTasks(userID string) ([]*model.TranscriptionTask, error) {
	return s.repo.ListTasks(userID)
}

func (s *TranscribeService) DeleteTranscriptionTask(userID, taskID string) error {
	return s.repo.DeleteTask(userID, taskID)
}

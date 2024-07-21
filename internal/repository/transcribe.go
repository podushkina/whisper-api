// internal/repository/transcribe.go
package repository

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8" // TODO: уже вышла сервия v9 переписать под нее
	"whisper-api/internal/model"
)

var ErrUnauthorized = errors.New("unauthorized access")

type TranscribeRepository struct {
	redis *redis.Client
}

func NewTranscribeRepository(redis *redis.Client) *TranscribeRepository {
	return &TranscribeRepository{redis: redis}
}

func (r *TranscribeRepository) SaveTask(task *model.TranscriptionTask) error {
	ctx := context.Background()
	taskJSON, err := json.Marshal(task)
	if err != nil {
		return err
	}
	return r.redis.Set(ctx, "task:"+task.ID, taskJSON, 0).Err()
}

func (r *TranscribeRepository) GetTask(userID, taskID string) (*model.TranscriptionTask, error) {
	ctx := context.Background()
	taskJSON, err := r.redis.Get(ctx, "task:"+taskID).Bytes()
	if err != nil {
		return nil, err
	}

	var task model.TranscriptionTask
	err = json.Unmarshal(taskJSON, &task)
	if err != nil {
		return nil, err
	}

	if task.UserID != userID {
		return nil, ErrUnauthorized
	}

	return &task, nil
}

func (r *TranscribeRepository) UpdateTask(task *model.TranscriptionTask) error {
	return r.SaveTask(task)
}

func (r *TranscribeRepository) ListTasks(userID string) ([]*model.TranscriptionTask, error) {
	ctx := context.Background()
	keys, err := r.redis.Keys(ctx, "task:*").Result()
	if err != nil {
		return nil, err
	}

	var tasks []*model.TranscriptionTask
	for _, key := range keys {
		taskJSON, err := r.redis.Get(ctx, key).Bytes()
		if err != nil {
			continue
		}

		var task model.TranscriptionTask
		err = json.Unmarshal(taskJSON, &task)
		if err != nil {
			continue
		}

		if task.UserID == userID {
			tasks = append(tasks, &task)
		}
	}

	return tasks, nil
}

func (r *TranscribeRepository) DeleteTask(userID, taskID string) error {
	ctx := context.Background()
	task, err := r.GetTask(userID, taskID)
	if err != nil {
		return err
	}

	return r.redis.Del(ctx, "task:"+task.ID).Err()
}

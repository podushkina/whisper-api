package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"whisper-api/internal/service"
	"whisper-api/pkg/whisper"
)

type TranscribeHandler struct {
	transcribeService *service.TranscribeService
}

func NewTranscribeHandler(transcribeService *service.TranscribeService) *TranscribeHandler {
	return &TranscribeHandler{transcribeService: transcribeService}
}

func (h *TranscribeHandler) CreateTranscription(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("audio")
	if err != nil {
		http.Error(w, "Invalid audio file", http.StatusBadRequest)
		return
	}
	defer file.Close() // TODO: написать хендерл для ошибки

	userID := r.Context().Value("user_id").(string)

	language := r.FormValue("language")
	outputFormat := r.FormValue("output_format")
	model := r.FormValue("model")
	temperature, _ := strconv.ParseFloat(r.FormValue("temperature"), 64)
	threads, _ := strconv.Atoi(r.FormValue("threads"))

	options := whisper.TranscriptionOptions{
		Language:     language,
		OutputFormat: outputFormat,
		Model:        model,
		Temperature:  temperature,
		Threads:      threads,
	}

	task, err := h.transcribeService.CreateTranscriptionTask(userID, file, header.Filename, options)
	if err != nil {
		http.Error(w, "Failed to create transcription task", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(task) // TODO: написать хендерл для ошибки
}

func (h *TranscribeHandler) GetTranscription(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	userID := r.Context().Value("user_id").(string)

	task, err := h.transcribeService.GetTranscriptionTask(userID, taskID)
	if err != nil {
		http.Error(w, "Transcription not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(task) // TODO: написать хендерл для ошибки
}

func (h *TranscribeHandler) ListTranscriptions(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	tasks, err := h.transcribeService.ListTranscriptionTasks(userID)
	if err != nil {
		http.Error(w, "Failed to list transcriptions", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tasks) // TODO: написать хендерл для ошибки
}

func (h *TranscribeHandler) DeleteTranscription(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	userID := r.Context().Value("user_id").(string)

	err := h.transcribeService.DeleteTranscriptionTask(userID, taskID)
	if err != nil {
		http.Error(w, "Failed to delete transcription", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// internal/handler/usage.go
package handler

import (
	"encoding/json"
	"net/http"

	"whisper-api/internal/service"
)

type UsageHandler struct {
	usageService *service.UsageService
}

func NewUsageHandler(usageService *service.UsageService) *UsageHandler {
	return &UsageHandler{usageService: usageService}
}

func (h *UsageHandler) GetUsage(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	usage, err := h.usageService.GetUsage(userID)
	if err != nil {
		http.Error(w, "Failed to get usage statistics", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(usage)
}

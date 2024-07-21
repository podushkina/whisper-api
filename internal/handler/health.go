// internal/handler/health.go
package handler

import (
	"encoding/json"
	"net/http"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	status := map[string]string{
		"status":  "OK",
		"version": "1.0.0",
	}
	json.NewEncoder(w).Encode(status) // TODO: написать хендерл для ошибкиф
}

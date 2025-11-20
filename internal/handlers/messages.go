package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/govnocods/RedChat/internal/logger"
)

func (h *Handlers) GetMessages(w http.ResponseWriter, r *http.Request) {
	messages, err := h.MessageService.GetMessage()
	if err != nil {
		logger.WithError(err).Error("Failed to get messages")
		http.Error(w, "Failed to get messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(messages); err != nil {
		logger.WithError(err).Error("Failed to encode messages response")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
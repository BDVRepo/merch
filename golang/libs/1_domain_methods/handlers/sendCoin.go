package handlers

import (
	"encoding/json"
	"net/http"
)

type SendCoinRequest struct {
	ToUserID int `json:"to_user_id"`
	Amount   int `json:"amount"`
}

// SendCoinHandler — отправка монет другому пользователю
func SendCoinHandler(w http.ResponseWriter, r *http.Request) {
	var req SendCoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Логика отправки монет
	w.WriteHeader(http.StatusOK)
}

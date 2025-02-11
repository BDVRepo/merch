package handlers

import (
	"net/http"
)

// BuyItemHandler — покупка предмета за монеты
func BuyItemHandler(w http.ResponseWriter, r *http.Request) {
	// item := chi.URLParam(r, "item")
	// Логика покупки
	w.WriteHeader(http.StatusOK)
}

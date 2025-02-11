package handlers

import (
	"encoding/json"
	"net/http"
)

type InfoResponse struct {
	Coins        int           `json:"coins"`
	Inventory    []string      `json:"inventory"`
	Transactions []Transaction `json:"transactions"`
}
type Transaction struct {
	ID     int    `json:"id"`
	Amount int    `json:"amount"`
	Type   string `json:"type"`
}

// InfoHandler — получить информацию о монетах, инвентаре и транзакциях
func InfoHandler(w http.ResponseWriter, r *http.Request) {
	// userID := r.Context().Value("user_id").(int) // предполагаем, что user_id уже в контексте
	// Логика получения информации
	response := InfoResponse{
		Coins:        100, // заглушка
		Inventory:    []string{"item1", "item2"},
		Transactions: []Transaction{},
	}
	json.NewEncoder(w).Encode(response)
}

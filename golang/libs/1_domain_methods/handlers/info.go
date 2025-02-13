package handlers

import (
	"bdv-avito-merch/libs/2_generated_models/model"
	"bdv-avito-merch/libs/4_common/smart_context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type InfoResponse struct {
	Coins        int32        `json:"coins"`
	Inventory    []Item       `json:"inventory"`
	CoinsHistory Transactions `json:"coins_history"`
}
type Item struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}
type Transactions struct {
	Sent     []Sent     `json:"sent"`
	Received []Received `json:"received"`
}
type Sent struct {
	ToUser string `json:"to_user"`
	Amount int32  `json:"amount"`
}
type Received struct {
	FromUser string `json:"from_user"`
	Amount   int32  `json:"amount"`
}

// InfoHandler — получить информацию о монетах, инвентаре и транзакциях
func InfoHandler(logger smart_context.ISmartContext, w http.ResponseWriter, r *http.Request) {
	fields := logger.GetDataFields()

	userID, ok := fields["UserID"].(string)
	if !ok {
		http.Error(w, "Не найден user_id в сессии", http.StatusUnauthorized)
		return
	}

	var user model.DocUser
	if err := logger.GetDB().First(&user, "user_id = ?", userID).Error; err != nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	// Получаем баланс пользователя
	coins := user.Balance

	// Запрос к инвентарю
	var inventoryData []struct {
		MerchCode string `gorm:"column:merch_code"`
		Count     int    `gorm:"column:count"`
	}
	if err := logger.GetDB().Table("doc_user_merchs").
		Select("merch_code, COUNT(merch_code) as count").
		Where("root_id = ?", user.ID).
		Group("merch_code").
		Scan(&inventoryData).Error; err != nil {
		http.Error(w, "Ошибка получения инвентаря", http.StatusInternalServerError)
		return
	}

	// Преобразуем в нужную структуру (избегаем nil)
	inventory := make([]Item, 0)
	for _, item := range inventoryData {
		inventory = append(inventory, Item{Name: item.MerchCode, Quantity: item.Count})
	}

	// Загружаем транзакции пользователя + получаем имя получателя/отправителя
	var transactions []struct {
		ID         uuid.UUID  `gorm:"column:id"`
		SenderID   uuid.UUID  `gorm:"column:sender_id"`
		ReceiverID *uuid.UUID `gorm:"column:receiver_id"`
		Amount     int32      `gorm:"column:amount"`
		CreatedAt  time.Time  `gorm:"column:created_at"`
		ToName     *string    `gorm:"column:to_name"`
		FromName   *string    `gorm:"column:from_name"`
	}
	if err := logger.GetDB().Raw(`
		SELECT t.id, t.sender_id, t.receiver_id, t.amount, t.created_at,
			   su.name AS from_name, ru.name AS to_name
		FROM doc_transactions t
		LEFT JOIN doc_users su ON t.sender_id = su.id
		LEFT JOIN doc_users ru ON t.receiver_id = ru.id
		WHERE t.sender_id = ? OR t.receiver_id = ?`, user.ID, user.ID).
		Scan(&transactions).Error; err != nil {
		http.Error(w, "Ошибка получения транзакций", http.StatusInternalServerError)
		return
	}

	// Разбиваем транзакции на отправленные и полученные
	sentTransactions := make([]Sent, 0)
	receivedTransactions := make([]Received, 0)

	for _, t := range transactions {
		if t.ReceiverID != nil && t.SenderID.String() == *user.ID {
			sentTransactions = append(sentTransactions, Sent{
				ToUser: *t.ToName,
				Amount: t.Amount,
			})
		} else if t.ReceiverID != nil && t.ReceiverID.String() == *user.ID {
			receivedTransactions = append(receivedTransactions, Received{
				FromUser: *t.FromName,
				Amount:   t.Amount,
			})
		}
	}

	// Формируем ответ (избегаем nil)
	response := InfoResponse{
		Coins:     coins,
		Inventory: inventory,
		CoinsHistory: Transactions{
			Sent:     sentTransactions,
			Received: receivedTransactions,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

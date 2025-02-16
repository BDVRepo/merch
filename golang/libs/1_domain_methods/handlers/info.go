package handlers

import (
	"bdv-avito-merch/libs/2_generated_models/model"
	"bdv-avito-merch/libs/4_common/smart_context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type InfoResponse struct {
	Name         string       `json:"name"`
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
type InfoRequest struct {
	logger smart_context.ISmartContext
	UserID string
}
type InfoQuery struct {
	r            InfoRequest
	responseChan chan InfoResponseData
}
type InfoResponseData struct {
	data   *InfoResponse
	err    error
	status int
}

func info(logger smart_context.ISmartContext, userID string) (*InfoResponse, error, int) {
	var user model.DocUser
	if err := logger.GetDB().First(&user, "user_id = ?", userID).Error; err != nil {
		return nil, fmt.Errorf("пользователь не найден"), http.StatusNotFound
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
		return nil, fmt.Errorf("ошибка получения инвентаря"), http.StatusInternalServerError
	}

	inventory := make([]Item, 0, len(inventoryData))
	for _, item := range inventoryData {
		inventory = append(inventory, Item{Name: item.MerchCode, Quantity: item.Count})
	}

	// Загружаем транзакции пользователя
	var transactions []struct {
		SenderID   uuid.UUID  `gorm:"column:sender_id"`
		ReceiverID *uuid.UUID `gorm:"column:receiver_id"`
		Amount     int32      `gorm:"column:amount"`
		ToName     *string    `gorm:"column:to_name"`
		FromName   *string    `gorm:"column:from_name"`
	}
	if err := logger.GetDB().Raw(`
		SELECT t.sender_id, t.receiver_id, t.amount, su.name AS from_name, ru.name AS to_name
		FROM doc_transactions t
		LEFT JOIN doc_users su ON t.sender_id = su.id
		LEFT JOIN doc_users ru ON t.receiver_id = ru.id
		WHERE t.sender_id = ? OR t.receiver_id = ?`, user.ID, user.ID).
		Scan(&transactions).Error; err != nil {
		return nil, fmt.Errorf("ошибка получения транзакций"), http.StatusInternalServerError
	}

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

	return &InfoResponse{
		Name:      user.Name,
		Coins:     coins,
		Inventory: inventory,
		CoinsHistory: Transactions{
			Sent:     sentTransactions,
			Received: receivedTransactions,
		},
	}, nil, http.StatusOK
}

// InfoHandler — получить информацию о монетах, инвентаре и транзакциях
func InfoHandler(logger smart_context.ISmartContext, w http.ResponseWriter, r *http.Request) {
	fields := logger.GetDataFields()
	userID, ok := fields["UserID"].(string)
	if !ok {
		http.Error(w, `{"errors": "Не авторизован"}`, http.StatusUnauthorized)
		return
	}

	responseChan := make(chan InfoResponseData)
	query := InfoQuery{
		r: InfoRequest{
			logger: logger,
			UserID: userID,
		},
		responseChan: responseChan,
	}

	// Отправляем запрос в воркер
	handlersRequests <- query
	response := <-responseChan

	if response.err != nil {
		http.Error(w, response.err.Error(), response.status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response.data); err != nil {
		http.Error(w, "Ошибка кодирования JSON", http.StatusInternalServerError)
		return
	}
}

package handlers

import (
	"bdv-avito-merch/libs/2_generated_models/model"
	"bdv-avito-merch/libs/4_common/smart_context"
	"encoding/json"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

// Структура для отправки монет
type SendCoinRequest struct {
	logger     smart_context.ISmartContext
	ToUsername string `json:"to_username"`
	Amount     int32  `json:"amount"`
}

// Функция для обработки транзакции
func sendCoinTransaction(req SendCoinRequest) error {

	return req.logger.GetDB().Transaction(func(tx *gorm.DB) error {
		var sender model.DocUser
		var receiver model.DocUser

		fields := req.logger.GetDataFields()
		sender_id, ok := fields["UserID"].(string)
		if !ok {
			return fmt.Errorf("Не нашел отправителя в сессии")
		}

		// 1. Получаем данные отправителя и получателя
		if err := tx.First(&sender).Where("user_id = ?", sender_id).Error; err != nil {
			return fmt.Errorf("Неизвестный отправитель: %w", err)
		}

		// 2. Получаем получателя по username с использованием Find()
		var receivers []model.DocUser
		if err := tx.Where("name = ?", req.ToUsername).Find(&receivers).Error; err != nil || len(receivers) == 0 {
			return fmt.Errorf("Ошибка при получении получателя: %w", err)
		}
		receiver = receivers[0]

		// 2. Проверка на достаточное количество монет у отправителя
		if sender.Balance < req.Amount {
			return fmt.Errorf("Недостаточно баланса отправителя")
		}

		// 3. Обновление баланса отправителя
		if err := tx.Model(&sender).Update("balance", sender.Balance-req.Amount).Error; err != nil {
			return fmt.Errorf("Не удалось обновить баланс отправителя: %w", err)
		}

		// 4. Обновление баланса получателя
		if err := tx.Model(&receiver).Update("balance", receiver.Balance+req.Amount).Error; err != nil {
			return fmt.Errorf("Не удалось обновить баланс получателя: %w", err)
		}

		// 5. Запись транзакции
		transaction := model.DocTransaction{
			SenderID:      *sender.ID,
			ReceiverID:    receiver.ID,
			OperationCode: "SEND",
			Amount:        req.Amount,
		}

		if err := tx.Create(&transaction).Error; err != nil {
			return fmt.Errorf("Не удалось записать транзакцию отправки монет: %w", err)
		}

		return nil
	})
}

// SendCoinHandler — отправка монет другому пользователю
func SendCoinHandler(logger smart_context.ISmartContext, w http.ResponseWriter, r *http.Request) {
	var req SendCoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	req.logger = logger

	// Ставим задачу в очередь
	balanceRequests <- req

	// TODO:Ответ об успешной отправке
	w.WriteHeader(http.StatusOK)
}

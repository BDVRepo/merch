package handlers

import (
	"bdv-avito-merch/libs/1_domain_methods/helpers"
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

type SendItemQuery struct {
	r            SendCoinRequest
	w            http.ResponseWriter
	responseChan chan Response
}

// Функция для отправки монет
func sendCoinTransaction(req SendCoinRequest) (err error, status int) {
	status = http.StatusOK
	err = req.logger.GetDB().Transaction(func(tx *gorm.DB) error {
		var sender model.DocUser
		var receiver model.DocUser

		fields := req.logger.GetDataFields()
		sender_id, ok := fields["UserID"].(string)
		if !ok {
			status = http.StatusBadRequest
			return fmt.Errorf("Не нашел отправителя в сессии")
		}

		// 1. Получаем данные отправителя и получателя
		if err := tx.First(&sender, "user_id = ?", sender_id).Error; err != nil {
			status = http.StatusNotFound
			return fmt.Errorf("Неизвестный отправитель: %w", err)
		}

		// 2. Получаем получателя по username с использованием Find()
		if err := tx.First(&receiver, "name = ?", req.ToUsername).Error; err != nil {
			status = http.StatusNotFound
			return fmt.Errorf("Ошибка при получении получателя: %w", err)
		}

		if *sender.ID == *receiver.ID {
			status = http.StatusBadRequest
			return fmt.Errorf("Нельзя отправить монеты самому себе")
		}

		// 2. Проверка на достаточное количество монет у отправителя
		if sender.Balance < req.Amount {
			status = http.StatusPaymentRequired
			return fmt.Errorf("Недостаточно баланса отправителя")
		}

		// 3. Обновление баланса отправителя
		if err := tx.Model(&sender).Update("balance", sender.Balance-req.Amount).Error; err != nil {
			status = http.StatusInternalServerError
			return fmt.Errorf("Не удалось обновить баланс отправителя: %w", err)
		}

		// 4. Обновление баланса получателя
		if err := tx.Model(&receiver).Update("balance", receiver.Balance+req.Amount).Error; err != nil {
			status = http.StatusInternalServerError
			return fmt.Errorf("Не удалось обновить баланс получателя: %w", err)
		}

		// 5. Запись транзакции
		transaction := model.DocTransaction{
			ID:         helpers.GenerateUUID(),
			SenderID:   *sender.ID,
			ReceiverID: receiver.ID,
			Amount:     req.Amount,
		}

		if err := tx.Create(&transaction).Error; err != nil {
			return fmt.Errorf("Не удалось записать транзакцию отправки монет: %w", err)
		}

		return nil
	})
	return err, status
}

// SendCoinHandler — отправка монет другому пользователю
func SendCoinHandler(logger smart_context.ISmartContext, w http.ResponseWriter, r *http.Request) {
	responseChan := make(chan Response)
	var req SendCoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	req.logger = logger

	query := SendItemQuery{
		r:            req,
		w:            w,
		responseChan: responseChan,
	}

	// Ставим задачу в очередь
	handlersRequests <- query

	response := <-responseChan

	// Возвращаем HTTP-ответ
	if response.err != nil {
		http.Error(w, response.err.Error(), response.status)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Токены успешно переданы"))
}

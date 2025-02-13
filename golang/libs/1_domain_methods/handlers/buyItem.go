package handlers

import (
	"bdv-avito-merch/libs/1_domain_methods/helpers"
	"bdv-avito-merch/libs/2_generated_models/model"
	"bdv-avito-merch/libs/4_common/smart_context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// Структура для покупки товара
type BuyItemRequest struct {
	logger    smart_context.ISmartContext
	MerchName string `json:"merch_name"`
}
type BuyItemQuery struct {
	r            BuyItemRequest
	w            http.ResponseWriter
	responseChan chan Response
}
type Response struct {
	err    error
	status int
}

// Функция для обработки покупки товара
func buyItemTransaction(req BuyItemRequest) (err error, status int) {
	status = http.StatusOK
	err = req.logger.GetDB().Transaction(func(tx *gorm.DB) error {
		if req.MerchName == "" {
			status = http.StatusBadRequest
			return fmt.Errorf("Название товара не может быть пустым")
		}

		// Получаем данные покупателя
		fields := req.logger.GetDataFields()
		sender_id, ok := fields["UserID"].(string)
		if !ok {
			status = http.StatusBadRequest
			return fmt.Errorf("Не нашел покупателя в сессии")
		}
		var buyer model.DocUser
		if err := tx.First(&buyer, "user_id = ?", sender_id).Error; err != nil {
			status = http.StatusNotFound
			return fmt.Errorf("Неизвестный покупатель: %w", err)
		}
		// Получаем товар для покупки
		var merch model.DocMerch
		if err := tx.First(&merch, "code = ?", req.MerchName).Error; err != nil {
			status = http.StatusNotFound
			return fmt.Errorf("Неизвестный товар: %w", err)
		}

		if buyer.Balance < merch.Price {
			status = http.StatusPaymentRequired
			return fmt.Errorf("Недостаточно баланса для покупки")
		}

		// Обновление баланса покупателя
		if err := tx.Model(&buyer).Update("balance", buyer.Balance-merch.Price).Error; err != nil {
			status = http.StatusInternalServerError
			return fmt.Errorf("Не удалось обновить баланс покупателя: %w", err)
		}

		// Запись мерча покупателю
		docUserMerch := model.DocUserMerch{
			ID:        helpers.GenerateUUID(),
			RootID:    *buyer.ID,
			MerchCode: merch.Code,
		}

		if err := tx.Create(&docUserMerch).Error; err != nil {
			status = http.StatusInternalServerError
			return fmt.Errorf("Не удалось записать покупку мерча покупателю: %w", err)
		}
		return nil
	})
	return err, status
}

// BuyItemHandler — покупка предмета за монеты
func BuyItemHandler(logger smart_context.ISmartContext, w http.ResponseWriter, r *http.Request) {

	merchName := chi.URLParam(r, "item")

	responseChan := make(chan Response)
	req := BuyItemRequest{
		logger:    logger,
		MerchName: merchName,
	}
	query := BuyItemQuery{
		r:            req,
		w:            w,
		responseChan: responseChan,
	}

	// Ставим задачу в очередь и ждем результат
	balanceRequests <- query

	response := <-responseChan

	// Возвращаем HTTP-ответ
	if response.err != nil {
		http.Error(w, response.err.Error(), response.status)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Мерч успешно куплен"))
}

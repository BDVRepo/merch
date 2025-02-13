package handlers

import (
	"bdv-avito-merch/libs/2_generated_models/model"
	"bdv-avito-merch/libs/4_common/smart_context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// Структура для покупки товара
type BuyItemRequest struct {
	logger    smart_context.ISmartContext
	MerchName string `json:"merch_name"`
}

// Функция для обработки покупки товара
func buyItemTransaction(req BuyItemRequest) error {
	return req.logger.GetDB().Transaction(func(tx *gorm.DB) error {
		var buyer model.DocUser
		var merch model.DocMerch

		// Получаем данные покупателя
		fields := req.logger.GetDataFields()
		sender_id, ok := fields["UserID"].(string)
		if !ok {
			return fmt.Errorf("Не нашел покупателя в сессии")
		}

		if err := tx.First(&buyer).Where("user_id = ?", sender_id).Error; err != nil {
			return fmt.Errorf("Неизвестный покупатель: %w", err)
		}

		// Получаем товар для покупки
		if err := tx.First(&merch).Where("name = ?", req.MerchName).Error; err != nil {
			return fmt.Errorf("Неизвестный товар: %w", err)
		}

		if buyer.Balance < merch.Price {
			return fmt.Errorf("Недостаточно баланса для покупки")
		}

		// Обновление баланса покупателя
		if err := tx.Model(&buyer).Update("balance", buyer.Balance-merch.Price).Error; err != nil {
			return fmt.Errorf("Не удалось обновить баланс покупателя: %w", err)
		}

		// Запись мерча покупателю
		docUserMerch := model.DocUserMerch{
			RootID:    *buyer.ID,
			MerchCode: merch.Code,
		}

		if err := tx.Create(&docUserMerch).Error; err != nil {
			return fmt.Errorf("Не удалось записать покупку мерча покупателю: %w", err)
		}
		return nil
	})
}

// BuyItemHandler — покупка предмета за монеты
func BuyItemHandler(logger smart_context.ISmartContext, w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	merchName := vars["item"]

	req := BuyItemRequest{
		logger:    logger,
		MerchName: merchName,
	}

	// Ставим задачу в очередь
	balanceRequests <- req

	// TODO:Ответ об успешной покупке
	w.WriteHeader(http.StatusOK)
}

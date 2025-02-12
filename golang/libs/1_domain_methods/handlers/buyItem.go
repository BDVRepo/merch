package handlers

import (
	"bdv-avito-merch/libs/4_common/smart_context"
	"net/http"
)

// BuyItemHandler — покупка предмета за монеты
func BuyItemHandler(logger smart_context.ISmartContext, w http.ResponseWriter, r *http.Request) {
	// item := chi.URLParam(r, "item")
	// Логика покупки
	w.WriteHeader(http.StatusOK)
}

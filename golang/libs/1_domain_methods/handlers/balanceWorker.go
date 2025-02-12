package handlers

import (
	"log"
)

// Канал для работы с запросами
var balanceRequests = make(chan interface{}, 1000) // Канал для различных задач

// Обработчик очереди задач
func BalanceWorker() {
	for {
		select {
		case req := <-balanceRequests:
			switch r := req.(type) {
			case SendCoinRequest:
				err := sendCoinTransaction(r)
				if err != nil {
					log.Println("Error processing send coin:", err)
				}
			case BuyItemRequest:
				err := buyItemTransaction(r)
				if err != nil {
					log.Println("Error processing buy item:", err)
				}
			}
		}
	}
}

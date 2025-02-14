package handlers

import "bdv-avito-merch/libs/4_common/safe_go"

// Канал для работы с запросами
var balanceRequests = make(chan interface{}, 100000) // Канал для различных задач

// Обработчик очереди задач
func BalanceWorker() {
	for {
		select {
		case query := <-balanceRequests:
			switch q := query.(type) {
			case SendItemQuery:
				safe_go.SafeGo(q.r.logger, func() {
					err, status := sendCoinTransaction(q.r)

					q.responseChan <- Response{
						err:    err,
						status: status,
					}
				})

			case BuyItemQuery:
				safe_go.SafeGo(q.r.logger, func() {
					err, status := buyItemTransaction(q.r)

					q.responseChan <- Response{
						err:    err,
						status: status,
					}
				})
			}
		}
	}
}

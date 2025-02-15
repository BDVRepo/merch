package handlers

import (
	"bdv-avito-merch/libs/4_common/safe_go"
)

// Канал для работы с запросами
var handlersRequests = make(chan interface{})

// Обработчик очереди задач
func HandlersWorker() {
	for {
		select {
		case query := <-handlersRequests:
			switch q := query.(type) {
			case LoginQuery:
				safe_go.SafeGo(q.r.logger, func() {
					token, status := login(q.r.logger, q.r.Login, q.r.Password)

					q.responseChan <- LoginResponse{
						token:  token,
						status: status,
					}
				})

			case InfoQuery:
				safe_go.SafeGo(q.r.logger, func() {
					data, err, status := info(q.r.logger, q.r.UserID)

					q.responseChan <- InfoResponseData{
						data:   data,
						err:    err,
						status: status,
					}
				})

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

package handlers

import (
	"bdv-avito-merch/libs/2_generated_models/model"
	"bdv-avito-merch/libs/4_common/smart_context"
	"net/http"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUnit_InfoFunction(t *testing.T) {

	// Открываем базу данных SQLite в памяти
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Не удалось открыть базу данных: %v", err)
	}

	// Мигрируем необходимые схемы
	if err := db.AutoMigrate(
		&model.AuthUser{},
		&model.DocUser{},
		&model.DocMerch{},
		&model.DocUserMerch{},
		&model.DocTransaction{},
	); err != nil {
		t.Fatalf("Не удалось выполнить миграцию базы данных: %v", err)
	}

	db = db.Begin()
	defer db.Rollback()

	logger := smart_context.NewSmartContext().WithDB(db)

	// 1. Создаём двух пользователей через функцию login
	_, statusSender := login(logger, "sender", "password")
	if statusSender != http.StatusOK {
		t.Fatalf("Не удалось создать/авторизовать пользователя sender, статус: %d", statusSender)
	}
	_, statusReceiver := login(logger, "receiver", "password")
	if statusReceiver != http.StatusOK {
		t.Fatalf("Не удалось создать/авторизовать пользователя receiver, статус: %d", statusReceiver)
	}

	// Извлекаем AuthUser по логину для последующего доступа к DocUser
	var authSender model.AuthUser
	if err := db.First(&authSender, "login = ?", "sender").Error; err != nil {
		t.Fatalf("Не удалось найти AuthUser sender: %v", err)
	}
	var authReceiver model.AuthUser
	if err := db.First(&authReceiver, "login = ?", "receiver").Error; err != nil {
		t.Fatalf("Не удалось найти AuthUser receiver: %v", err)
	}

	// Извлекаем соответствующих DocUser
	var docSender model.DocUser
	if err := db.First(&docSender, "user_id = ?", authSender.ID).Error; err != nil {
		t.Fatalf("Не удалось найти DocUser sender: %v", err)
	}
	var docReceiver model.DocUser
	if err := db.First(&docReceiver, "user_id = ?", authReceiver.ID).Error; err != nil {
		t.Fatalf("Не удалось найти DocUser receiver: %v", err)
	}

	// Задаём начальные балансы и имена для наглядности
	if err := db.Model(&docSender).Updates(map[string]interface{}{"balance": 1000, "name": "Sender"}).Error; err != nil {
		t.Fatalf("Не удалось обновить баланс/имя sender: %v", err)
	}
	if err := db.Model(&docReceiver).Updates(map[string]interface{}{"balance": 500, "name": "Receiver"}).Error; err != nil {
		t.Fatalf("Не удалось обновить баланс/имя receiver: %v", err)
	}

	// 3. Пользователь sender покупает мерч "t-shirt = 80" через buyItemTransaction
	logger = logger.WithField("UserID", *authSender.ID)
	buyReq := BuyItemRequest{
		logger:    logger,
		MerchName: "t-shirt",
	}
	if err, _ := buyItemTransaction(buyReq); err != nil {
		t.Fatalf("buyItemTransaction завершилась ошибкой: %v", err)
	}
	// После покупки баланс sender становится: 1000 - 80 = 920

	// 4. Отправка монет: sender переводит 200 монет receiver
	// Баланс sender становится: 920 - 200 = 720, а receiver: 500 + 200 = 700
	sendReq1 := SendCoinRequest{
		logger:     logger,
		ToUsername: "Receiver", // Имя получателя должно совпадать с docReceiver.name
		Amount:     200,
	}
	if err, status := sendCoinTransaction(sendReq1); err != nil || status != 200 {
		t.Fatalf("sendCoinTransaction от sender завершилась ошибкой: %v, статус: %d", err, status)
	}

	// 5. Отправка монет: receiver переводит 150 монет sender
	// Баланс receiver становится: 700 - 150 = 550, а sender: 720 + 150 = 870
	logger = logger.WithField("UserID", *authReceiver.ID)
	sendReq2 := SendCoinRequest{
		logger:     logger,
		ToUsername: "Sender",
		Amount:     150,
	}
	if err, status := sendCoinTransaction(sendReq2); err != nil || status != 200 {
		t.Fatalf("sendCoinTransaction от receiver завершилась ошибкой: %v, статус: %d", err, status)
	}

	// 6. Получаем информацию о пользователе sender через функцию info
	infoResp, err, status := info(logger, docSender.UserID)
	if err != nil {
		t.Fatalf("Функция info завершилась ошибкой: %v", err)
	}
	// Проверяем статус 200
	if status != 200 {
		t.Errorf("Ожидался статус 200, получено %d", status)
	}
	// Проверяем итоговый баланс пользователя sender
	expectedBalance := int32(870) // 1000 - 20 (покупка) - 200 (отправка) + 150 (получено) = 930
	if infoResp.Coins != expectedBalance {
		t.Errorf("Ожидался баланс %d, получено %d", expectedBalance, infoResp.Coins)
	}

	// Проверяем, что в инвентаре присутствует товар "t-shirt" в количестве 1
	if len(infoResp.Inventory) != 1 {
		t.Errorf("Ожидался 1 элемент инвентаря, получено %d", len(infoResp.Inventory))
	} else {
		item := infoResp.Inventory[0]
		if item.Name != "t-shirt" {
			t.Errorf("Ожидалось, что товар будет 't-shirt', получено %s", item.Name)
		}
		if item.Quantity != 1 {
			t.Errorf("Ожидалось количество 1, получено %d", item.Quantity)
		}
	}

	// Проверяем историю транзакций
	// Для sender ожидается:
	// - Отправленная транзакция: получатель "Receiver", сумма 200
	// - Полученная транзакция: отправитель "Receiver", сумма 150
	if len(infoResp.CoinsHistory.Sent) != 1 {
		t.Errorf("Ожидалась 1 отправленная транзакция, получено %d", len(infoResp.CoinsHistory.Sent))
	} else {
		sentTx := infoResp.CoinsHistory.Sent[0]
		if sentTx.ToUser != "Receiver" {
			t.Errorf("Ожидалось, что отправлено пользователю 'Receiver', получено %s", sentTx.ToUser)
		}
		if sentTx.Amount != 200 {
			t.Errorf("Ожидалась сумма 200, получено %d", sentTx.Amount)
		}
	}
	if len(infoResp.CoinsHistory.Received) != 1 {
		t.Errorf("Ожидалась 1 полученная транзакция, получено %d", len(infoResp.CoinsHistory.Received))
	} else {
		receivedTx := infoResp.CoinsHistory.Received[0]
		if receivedTx.FromUser != "Receiver" {
			t.Errorf("Ожидалось, что получено от 'Receiver', получено %s", receivedTx.FromUser)
		}
		if receivedTx.Amount != 150 {
			t.Errorf("Ожидалась сумма 150, получено %d", receivedTx.Amount)
		}
	}
}

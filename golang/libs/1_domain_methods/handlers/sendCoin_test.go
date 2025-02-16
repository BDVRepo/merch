package handlers

import (
	"bdv-avito-merch/libs/1_domain_methods/helpers"
	"bdv-avito-merch/libs/2_generated_models/model"
	"bdv-avito-merch/libs/4_common/smart_context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUnit_SendCoinTransaction(t *testing.T) {
	// Открытие базы данных SQLite в памяти
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	// Миграция схемы
	if err := db.AutoMigrate(&model.AuthUser{}, &model.DocUser{}, &model.DocTransaction{}); err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// Начало транзакции
	db = db.Begin()
	defer db.Rollback() // Очищаем БД после теста

	// Создание логгера с контекстом базы данных
	logger := smart_context.NewSmartContext().WithDB(db)

	// Заполнение тестовыми данными
	senderAuth := model.AuthUser{Login: "sender", Password: "password", ID: helpers.GenerateUUID()}
	receiverAuth := model.AuthUser{Login: "receiver", Password: "password", ID: helpers.GenerateUUID()}
	if err := db.Create(&senderAuth).Error; err != nil {
		t.Fatalf("failed to create sender auth: %v", err)
	}
	if err := db.Create(&receiverAuth).Error; err != nil {
		t.Fatalf("failed to create receiver auth: %v", err)
	}

	sender := model.DocUser{UserID: *senderAuth.ID, Name: "Sender", Balance: 1000, ID: helpers.GenerateUUID()}
	receiver := model.DocUser{UserID: *receiverAuth.ID, Name: "Receiver", Balance: 500, ID: helpers.GenerateUUID()}
	if err := db.Create(&sender).Error; err != nil {
		t.Fatalf("failed to create sender: %v", err)
	}
	if err := db.Create(&receiver).Error; err != nil {
		t.Fatalf("failed to create receiver: %v", err)
	}

	logger = logger.WithField("UserID", *senderAuth.ID)

	tests := []struct {
		name       string
		toUsername string
		amount     int32
		wantErr    bool
		wantStatus int
		errMsg     string
	}{
		{
			name:       "Test valid transaction",
			toUsername: "Receiver",
			amount:     200,
			wantErr:    false,
			wantStatus: 200,
		},
		{
			name:       "Test insufficient balance",
			toUsername: "Receiver",
			amount:     2000,
			wantErr:    true,
			wantStatus: 402,
			errMsg:     "Недостаточно баланса отправителя",
		},
		{
			name:       "Test unknown receiver",
			toUsername: "UnknownUser",
			amount:     100,
			wantErr:    true,
			wantStatus: 404,
			errMsg:     "Ошибка при получении получателя",
		},
		{
			name:       "Test send coins to self",
			toUsername: "Sender", // отправка себе
			amount:     100,
			wantErr:    true,
			wantStatus: 400,
			errMsg:     "Нельзя отправить монеты самому себе",
		},
		{
			name:       "Test invalid request (missing amount)",
			toUsername: "Receiver",
			amount:     0,
			wantErr:    true,
			wantStatus: 400,
			errMsg:     "Неверные данные запроса: сумма должна быть больше 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := SendCoinRequest{
				logger:     logger,
				ToUsername: tt.toUsername,
				Amount:     tt.amount,
			}

			err, status := sendCoinTransaction(req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.wantStatus, status)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantStatus, status)

				var updatedSender model.DocUser
				var updatedReceiver model.DocUser

				if err := db.First(&updatedSender, "id = ?", sender.ID).Error; err != nil {
					t.Fatalf("failed to find updated sender: %v", err)
				}
				if err := db.First(&updatedReceiver, "id = ?", receiver.ID).Error; err != nil {
					t.Fatalf("failed to find updated receiver: %v", err)
				}

				assert.Equal(t, sender.Balance-tt.amount, updatedSender.Balance)
				assert.Equal(t, receiver.Balance+tt.amount, updatedReceiver.Balance)
			}
		})
	}
}

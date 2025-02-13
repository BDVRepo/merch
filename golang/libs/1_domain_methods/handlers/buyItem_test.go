package handlers

import (
	"bdv-avito-merch/libs/2_generated_models/model"
	"bdv-avito-merch/libs/4_common/smart_context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUnit_BuyItemTransaction(t *testing.T) {
	// Открытие базы данных SQLite в памяти
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	// Миграция схемы
	if err := db.AutoMigrate(
		&model.AuthUser{},
		&model.DocUser{},
		&model.DocMerch{},
		&model.DocUserMerch{},
	); err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// Начало транзакции
	db = db.Begin()
	defer db.Rollback() // Это очищает БД после завершения теста

	// Создание логгера с контекстом базы данных
	logger := smart_context.NewSmartContext().WithDB(db)

	// Заполнение тестовыми данными
	login(logger, "user123", "password123")
	buyer := model.DocUser{UserID: "user123"}
	merch := model.DocMerch{Code: "t-shirt", Price: 20}

	if err := db.Create(&merch).Error; err != nil {
		t.Fatalf("failed to create merch: %v", err)
	}

	tests := []struct {
		name      string
		merchName string
		userID    string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "Test valid purchase",
			merchName: "t-shirt",
			userID:    "user123",
			wantErr:   false,
		},
		{
			name:      "Test insufficient balance",
			merchName: "t-shirt",
			userID:    "user124", // Не существует в БД
			wantErr:   true,
			errMsg:    "Неизвестный покупатель",
		},
		{
			name:      "Test unknown merch",
			merchName: "UnknownItem", // Не существует в БД
			userID:    "user123",
			wantErr:   true,
			errMsg:    "Неизвестный товар",
		},
		{
			name:      "Test insufficient balance",
			merchName: "t-shirt",
			userID:    "user123",
			wantErr:   true,
			errMsg:    "Недостаточно баланса для покупки",
		},
	}

	// Прогон тестов
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Обновляем баланс покупателя для теста "Test insufficient balance"
			if tt.name == "Test insufficient balance" {
				// Если проверяем на недостаток баланса, уменьшаем баланс покупателя
				if err := db.Model(&buyer).Update("balance", 5.0).Error; err != nil {
					t.Fatalf("failed to update buyer's balance: %v", err)
				}
			}

			// Запрос на покупку товара
			req := BuyItemRequest{
				logger:    logger,
				MerchName: tt.merchName,
			}

			// Выполнение транзакции покупки
			err, _ := buyItemTransaction(req)

			// Проверка на ошибки
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)

				// Проверка, что баланс покупателя обновился правильно
				var updatedBuyer model.DocUser
				if err := db.First(&updatedBuyer, "user_id = ?", tt.userID).Error; err != nil {
					t.Fatalf("failed to find updated buyer: %v", err)
				}

				var updatedMerch model.DocUserMerch
				if err := db.First(&updatedMerch, "root_id = ? AND merch_code = ?", *updatedBuyer.ID, tt.merchName).Error; err != nil {
					t.Fatalf("failed to find purchased merch: %v", err)
				}

				// Проверка, что покупка была записана в DocUserMerch
				assert.Equal(t, tt.merchName, updatedMerch.MerchCode)
			}
		})
	}
}

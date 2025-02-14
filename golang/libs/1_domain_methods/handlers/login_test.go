package handlers

import (
	"bdv-avito-merch/libs/2_generated_models/model"
	"bdv-avito-merch/libs/4_common/env_vars"
	"bdv-avito-merch/libs/4_common/smart_context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Test_Login(t *testing.T) {
	env_vars.LoadEnvVars()
	// Создаём временную БД в памяти
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	// Миграция схемы БД
	if err := db.AutoMigrate(&model.AuthUser{}, &model.DocUser{}); err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	logger := smart_context.NewSmartContext().WithDB(db)
	password := "securepassword"

	tests := []struct {
		name       string
		login      string
		password   string
		wantErr    bool
		wantStatus int
		wantResp   string
	}{
		{
			name:       "Create user and login",
			login:      "new_user",
			password:   password,
			wantErr:    false,
			wantStatus: http.StatusOK,
			wantResp:   "eyJ",
		},
		{
			name:       "Login with correct password",
			login:      "new_user",
			password:   password,
			wantErr:    false,
			wantStatus: http.StatusOK,
			wantResp:   "eyJ",
		},
		{
			name:       "Login with incorrect password",
			login:      "new_user",
			password:   "wrongpassword",
			wantErr:    true,
			wantStatus: http.StatusUnauthorized,
			wantResp:   `{"errors": "Неверные учетные данные"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, status := login(logger, tt.login, tt.password)

			assert.Equal(t, tt.wantStatus, status)
			assert.Contains(t, resp, tt.wantResp)
		})
	}
}

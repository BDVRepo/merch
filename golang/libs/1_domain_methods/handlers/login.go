package handlers

import (
	"bdv-avito-merch/libs/1_domain_methods/helpers"
	"bdv-avito-merch/libs/2_generated_models/model"
	"bdv-avito-merch/libs/4_common/smart_context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type LoginRequest struct {
	logger   smart_context.ISmartContext
	Login    string
	Password string
}
type LoginQuery struct {
	r            LoginRequest
	responseChan chan LoginResponse
}
type LoginResponse struct {
	token  string
	status int
	err    error
}
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID string) (string, error) {
	jwtSecret, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		return "", errors.New("JWT_SECRET is not set")
	}
	expirationTime := time.Now().Add(2 * time.Hour)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
func login(logger smart_context.ISmartContext, login string, password string) (resp string, status int) {
	// Ищем пользователя
	var user model.AuthUser
	if err := logger.GetDB().First(&user, "login = ?", login).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Создаём нового пользователя
			hashedPassword, err := helpers.HashPassword(password)
			if err != nil {
				return `{"errors": "Ошибка хэширования пароля"}`, http.StatusInternalServerError
			}

			user = model.AuthUser{
				ID:       helpers.GenerateUUID(),
				Login:    login,
				Password: hashedPassword,
			}
			if err := logger.GetDB().Create(&user).Error; err != nil {
				return `{"errors": "Ошибка при создании пользователя"}`, http.StatusInternalServerError
			}

			// Создаём профиль
			docUser := model.DocUser{
				ID:      helpers.GenerateUUID(),
				UserID:  *user.ID,
				Name:    login,
				Balance: 1000,
			}
			if err := logger.GetDB().Create(&docUser).Error; err != nil {
				return `{"errors": "Ошибка при создании профиля"}`, http.StatusInternalServerError
			}
		} else {
			return `{"errors": "Ошибка базы данных"}`, http.StatusInternalServerError
		}
	}

	// Проверяем пароль
	if !helpers.CheckPassword(user.Password, password) {
		return `{"errors": "Неверные учетные данные"}`, http.StatusUnauthorized
	}

	// Генерируем JWT
	token, err := GenerateToken(*user.ID)
	if err != nil {
		return `{"errors": "Ошибка генерации токена"}`, http.StatusInternalServerError
	}
	return token, http.StatusOK
}
func LoginHandler(logger smart_context.ISmartContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestData struct {
			Login    string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, `{"errors": "Неверный запрос"}`, http.StatusBadRequest)
			return
		}

		responseChan := make(chan LoginResponse)
		query := LoginQuery{
			r: LoginRequest{
				logger:   logger,
				Login:    requestData.Login,
				Password: requestData.Password,
			},
			responseChan: responseChan,
		}

		// Ставим запрос в очередь
		handlersRequests <- query
		response := <-responseChan

		// Возвращаем HTTP-ответ
		if response.err != nil {
			http.Error(w, response.err.Error(), response.status)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{"token": response.token}); err != nil {
			http.Error(w, "Ошибка кодирования JSON", http.StatusInternalServerError)
			return
		}

	}
}

package handlers

import (
	"bdv-avito-merch/libs/1_domain_methods/helpers"
	"bdv-avito-merch/libs/2_generated_models/model"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

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

func LoginHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestData struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		login, password := requestData.Login, requestData.Password
		var user model.AuthUser

		// Ищем пользователя
		if err := db.Where("login = ?", login).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Создаём нового пользователя
				hashedPassword, err := helpers.HashPassword(password)
				if err != nil {
					http.Error(w, "Error hashing password", http.StatusInternalServerError)
					return
				}

				user = model.AuthUser{
					ID:       helpers.GenerateUUID(),
					Login:    login,
					Password: hashedPassword,
				}
				if err := db.Create(&user).Error; err != nil {
					http.Error(w, "Error creating user", http.StatusInternalServerError)
					return
				}

				// Создаём профиль
				infoUser := model.InfoUser{
					ID:      helpers.GenerateUUID(),
					UserID:  *user.ID,
					Name:    login,
					Balance: 1000,
				}
				if err := db.Create(&infoUser).Error; err != nil {
					http.Error(w, "Error creating profile", http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
		}

		// Проверяем пароль
		if !helpers.CheckPassword(user.Password, password) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Генерируем JWT
		token, err := GenerateToken(*user.ID)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		// Отправляем токен в теле ответа
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}

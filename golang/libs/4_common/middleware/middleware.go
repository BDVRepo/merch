package middleware

import (
	"bdv-avito-merch/libs/4_common/smart_context"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

// Claims - структура для хранения данных токена
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func WithSmartContext(
	externalSctx smart_context.ISmartContext,
	handler func(sctx smart_context.ISmartContext, w http.ResponseWriter, r *http.Request),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		jwtSecret, ok := os.LookupEnv("JWT_SECRET")
		if !ok {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return []byte(jwtSecret), nil
		})
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		authSctx := externalSctx.WithField("UserID", token.Claims.(*Claims).UserID)

		// Pass the updated SmartContext to the handler
		handler(authSctx, w, r)
	}
}

# Используем официальный образ Go
FROM golang:1.22 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod tidy

# Копируем весь проект
COPY . .

# Собираем приложение
RUN go build -o /build ./apps/backend-chi/main.go && go clean -cache -modcache

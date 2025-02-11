# Этап 1: Сборка
FROM golang:1.22 AS builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod tidy

# Копируем все файлы проекта
COPY . .

# Печатаем текущую директорию, чтобы убедиться, что копирование прошло корректно
RUN pwd
RUN ls -la

# Собираем приложение
RUN go build -o main ./apps/backend-chi/main.go

# Этап 2: Запуск
FROM golang:1.22

WORKDIR /root/

# Копируем скомпилированное приложение из builder
COPY --from=builder /app/main ./

EXPOSE 8080

# Запускаем приложение
CMD ["./main"]

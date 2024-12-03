# Этап сборки
FROM golang:1.23 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Сборка приложения
RUN go build -o app .

# Финальный образ
FROM debian:bullseye-slim

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем собранное приложение
COPY --from=builder /app/app .
COPY config.env .

# Указываем порт, который будет слушать приложение
EXPOSE 8080

# Запуск приложения
CMD ["./app"]
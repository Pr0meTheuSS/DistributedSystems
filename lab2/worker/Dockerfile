# Этап 1: сборка бинарника
FROM golang:1.24 AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum (если они есть) для установки зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем весь исходный код
COPY . .

# Сборка бинарника
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./internal/cmd

# Этап 2: минимальный образ для запуска
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем бинарник из предыдущего этапа
COPY --from=builder /app/app .

# Экспонируем порт (если нужно, например, 8080)
EXPOSE 9091

# Запускаем бинарник
ENTRYPOINT ["./app"]

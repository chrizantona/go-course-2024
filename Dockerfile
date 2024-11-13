# Стадия сборки
FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o storage ./cmd/main.go

# Финальный образ
FROM ubuntu:22.04

WORKDIR /root/

# Создаём директорию для данных
RUN mkdir -p /root/storage_data

# Копируем приложение
COPY --from=builder /app/storage .

RUN chmod +x /root/storage

EXPOSE 8090

CMD ["./storage"]

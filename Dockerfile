FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o storage ./cmd/main.go

FROM ubuntu:22.04

WORKDIR /root/


RUN mkdir -p /root/storage_data

COPY --from=builder /app/storage .

RUN chmod +x /root/storage

EXPOSE 8090

CMD ["./storage"]

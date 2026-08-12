# Шаг 1: Сборка бинарника (Builder)
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем бинарный файл под Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o /blog-api ./cmd/api/main.go

# Шаг 2: Легковесный финальный образ
FROM alpine:latest

WORKDIR /app

# Копируем только скомпилированный файл из builder
COPY --from=builder /blog-api /app/blog-api

EXPOSE 8080

CMD ["/app/blog-api"]
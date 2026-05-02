FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download || true

COPY . .

EXPOSE 8080

CMD ["go", "run", "./cmd/api/main.go"]
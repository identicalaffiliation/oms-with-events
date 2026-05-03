FROM golang:1.26.0 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o order-service ./cmd/api

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/order-service .

CMD ["./order-service", "-c", "config.yml"]
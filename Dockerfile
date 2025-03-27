# Stage 1: Build Go binary
FROM golang:1.24.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем бинарник из каталога cmd (где находится main.go)
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/app ./cmd

# Stage 2: Final image (Alpine)
FROM alpine:3.17

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Копируем собранный бинарник из builder
COPY --from=builder /go/bin/app /app/app

# КОПИРУЕМ ПАПКУ certs (где лежат cert.pem, key.pem)
COPY certs/ /app/certs/

USER appuser

# Ваше приложение слушает 8443 (по коду)
EXPOSE 8443

CMD ["/app/app"]

# Этап сборки
FROM golang:1.24.1-alpine AS builder

# Установка зависимостей для сборки (опционально, но полезно)
RUN apk add --no-cache git

WORKDIR /app


COPY go.mod ./

RUN go mod download


COPY . .


RUN CGO_ENABLED=0 GOOS=linux go build -v -o main ./cmd/app/main.go


FROM alpine:latest

WORKDIR /app


COPY --from=builder /app/main .


RUN chmod +x /app/main


EXPOSE 8080


ENTRYPOINT ["/app/main"]
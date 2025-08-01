FROM --platform=linux/arm64 golang:1.24.1-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=off

COPY order_service/go.mod order_service/go.sum ./
COPY shared ../shared

RUN go mod download

COPY order_service/cmd ./cmd
COPY order_service/internal ./internal

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o main ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add wget

COPY --from=builder /app/main /app/main

RUN chmod +x /app/main

WORKDIR /app

EXPOSE 8080

ENTRYPOINT ["/app/main"]
FROM golang:1.24.1-alpine AS builder

WORKDIR /app

COPY go.mod ./

COPY . .

RUN CGO_ENABLED=0 GOOS=darwin go build -v -o main ./cmd/app/main.go

RUN chmod +x /app/main

FROM alpine:latest

COPY --from=builder /app/main /app/main

WORKDIR /app

RUN chmod +x /app/main

EXPOSE 8080

ENTRYPOINT ["/app/main"]
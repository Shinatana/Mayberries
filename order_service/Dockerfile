FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=off

COPY infrastructure/go.mod infrastructure/go.sum ./
RUN go mod download

COPY infrastructure/cmd ./cmd
COPY infrastructure/internal ./internal
COPY infrastructure/pkg ./pkg

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:latest


COPY --from=builder /app/main /app/main

WORKDIR /app

EXPOSE 8080

ENTRYPOINT ["/app/main"]
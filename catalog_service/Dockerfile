FROM --platform=linux/arm64 golang:1.24.1-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

ENV GOPROXY=https://goproxy.io,direct
ENV GOSUMDB=off

COPY go.mod go.sum ./

RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o main ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add wget

COPY --from=builder /app/main /app/main

RUN chmod +x /app/main

WORKDIR /app

EXPOSE 8080

ENTRYPOINT ["/app/main"]
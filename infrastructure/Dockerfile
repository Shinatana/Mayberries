FROM golang:1.24.1-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=off

COPY go.mod ./

RUN go mod download

COPY cmd ./cmd

RUN CGO_ENABLED=0 GOOS=darwin go build -o main ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

ENTRYPOINT ["./main"]
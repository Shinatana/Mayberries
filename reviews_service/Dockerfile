FROM golang:1.24.1-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=off

COPY go.mod ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal


RUN CGO_ENABLED=0 GOOS=darwin go build -o main ./cmd/main.go

FROM alpine:latest


COPY --from=builder /app/main /app/main

WORKDIR /app

EXPOSE 8080

ENTRYPOINT ["/app/main"]
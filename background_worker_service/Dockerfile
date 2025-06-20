FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=off

COPY go.mod ./
RUN go mod download

COPY cmd ./cmd

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app/main.go

RUN chmod +x /app/main

FROM alpine:latest


COPY --from=builder /app/main /app/main

WORKDIR /app

RUN chmod +x /app/main

EXPOSE 8080

ENTRYPOINT ["/app/main"]
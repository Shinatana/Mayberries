FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

ENV Goproxy=https://goproxy.io,direct
ENV GOSUMDB=off

COPY catalog-service/go.mod catalog/go.sum ./

RUN go mod download

COPY catalog-service/cmd ./cmd
COPY catalog-service/internal ./internal
COPY catalog-service/pkg ./pkg

RUN CGO_ENABLE=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add wget

COPY --from=builder /app/main /app/main

WORKDIR /app

EXPOSE 8081

ENTRYPOINT ["/app/main"]








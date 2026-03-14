FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o brainless-waf cmd/brainless-waf/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/brainless-waf .
COPY config/config.yaml config/config.yaml

EXPOSE 80

CMD ["./brainless-waf"]

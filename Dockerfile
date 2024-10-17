# Keep this simple and clean
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod ./

COPY . .

RUN go build -o redis-cli-app

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/redis-cli-app .

CMD ["./redis-cli-app"]
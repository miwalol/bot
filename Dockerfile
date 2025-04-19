FROM golang:1.24 AS builder

WORKDIR /build

COPY . .

RUN go mod tidy && \
    go build -o bot .

FROM debian:stable-slim

WORKDIR /app

RUN apt update && apt install -y ca-certificates

COPY --from=builder /build/bot /app/bot

ENV GIN_MODE=release

CMD ["/app/bot"]
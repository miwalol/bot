FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY . .

RUN go mod tidy && \
    go build -o bot .

FROM alpine:latest

WORKDIR /app

RUN apk update && apk add --no-cache ca-certificates

COPY --from=builder /build/bot /app/bot

ENV GIN_MODE=release

CMD ["/app/bot"]
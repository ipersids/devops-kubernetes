# Build stage
FROM golang:alpine AS builder

WORKDIR /app
COPY . .

RUN go build -o todobackend .

# Runtime stage
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/todobackend ./todobackend

CMD ["./todobackend"]

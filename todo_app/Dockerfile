# Build stage
FROM golang:alpine AS builder

WORKDIR /app
COPY . .

RUN go build -o todoapp .

# Runtime stage
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/todoapp ./todoapp

CMD ["./todoapp"]

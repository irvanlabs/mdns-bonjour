# Stage 1: Build
FROM golang:alpine3.21 AS builder

RUN apk add --no-cache git build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build binary dari main.go
RUN go build -o server .

# Stage 2: Run
FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/server .

# Expose Bonjour & HTTP
EXPOSE 49221
EXPOSE 80

# Jalankan binary
CMD ["./server"]

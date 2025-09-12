# Start from the official Golang image
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o url-shortner .

# Use a minimal image for running
FROM alpine:latest

WORKDIR /app

# Copy the built binary and .env file
COPY --from=builder /app/url-shortner .

# Install 'grep' for extracting port from .env
RUN apk add --no-cache grep


EXPOSE 8001

# Run the application
CMD ["./url-shortner"]
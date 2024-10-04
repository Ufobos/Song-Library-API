# Build stage
FROM golang:1.23-alpine AS builder


WORKDIR /app

# Install git (required for fetching dependencies)
RUN apk update && apk add --no-cache git

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main ./cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy migrations and .env file
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/.env .env

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./main"]

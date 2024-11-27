FROM golang:alpine

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main cmd/main.go

# Create directories
RUN mkdir -p /app/logs /app/data

# Copy environment file
COPY .env.sample .env

# Set permissions
RUN chmod +x main

# Network configuration
ENV SERVER_HOST=0.0.0.0
ENV SERVER_PORT=8000

# Expose API port
EXPOSE 8000

# Define mount points for logs and data
VOLUME ["/app/logs", "/app/data"]

CMD ["./main"]
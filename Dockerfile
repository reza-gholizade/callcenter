# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Set Go toolchain
ENV GOTOOLCHAIN=local

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/main ./cmd/api

# Final stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Set working directory
WORKDIR /app

# Create bin directory
RUN mkdir -p /app/bin

# Copy the binary from builder
COPY --from=builder /app/bin/main /app/bin/

# Copy config files
COPY .env.example .env

# Expose port
EXPOSE 8080

# Run the application
CMD ["/app/bin/main"] 
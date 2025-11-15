# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Set Go toolchain to auto-download if needed for go.mod version
ENV GOTOOLCHAIN=auto

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies (will auto-download Go 1.24 if needed)
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 5000

# Run application
CMD ["./main"]


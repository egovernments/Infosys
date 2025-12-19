# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application (main.go is in cmd/server/)
RUN CGO_ENABLED=0 GOOS=linux go build -o property-enumeration-service ./cmd/server

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/property-enumeration-service .

# Set permissions
RUN chmod +x /app/property-enumeration-service

# Expose port 8080 (default property service port)
EXPOSE 8080

# Set the entry point
ENTRYPOINT ["./property-enumeration-service"]
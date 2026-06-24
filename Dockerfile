# Build stage
FROM golang:1.26.4-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# copy .env file if exists
RUN cp -r .env.example .env || true


# Build the application
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /home/appuser

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy migrations if needed at runtime
COPY --from=builder /app/internal/migrations ./internal/migrations

# Copy .env file if exists (optional - better to use env vars in production)
COPY --from=builder /app/.env .env* ./

# Change ownership to non-root user
RUN chown -R appuser:appuser /home/appuser

# Switch to non-root user
USER appuser

# Expose port (adjust based on your application)
EXPOSE 8080

# Run the application
CMD ["./main"]

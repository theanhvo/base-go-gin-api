# Multi-stage build for optimized production image
FROM golang:1.23-alpine AS builder

# Install build dependencies - cached layer
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy dependency files first for better caching
# This layer only rebuilds when go.mod/go.sum changes
COPY go.mod go.sum ./

# Download dependencies - cached unless go.mod/go.sum changes
RUN go mod download && go mod verify

# Copy source code last to maximize cache hits
COPY . .

# Build optimized binary with build cache
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -trimpath \
    -o main main.go

# Final production image
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata wget && \
    update-ca-certificates

# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Create app directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Set proper ownership
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port (configurable via K8s service)
EXPOSE 8080

# Health check endpoint for K8s probes
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider --timeout=5 http://localhost:8080/health || exit 1

# Run application
CMD ["./main"]

# Build stage
FROM golang:1.25.5-alpine3.23 AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o /app/noaas .

# Final stage - minimal alpine image
FROM alpine:3.23

WORKDIR /app

# Install ca-certificates for HTTPS and wget for health checks
RUN apk add --no-cache ca-certificates wget

# Copy binary from builder
COPY --from=builder /app/noaas .

# Copy data files (reasons)
COPY --from=builder /app/data ./data

# Expose port
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/app/noaas"]

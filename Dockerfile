FROM golang:1.21-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git make

# Copy only the necessary files for dependency download
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN make build

# Create final lightweight image
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy the binary from builder
COPY --from=builder /build/build/muraena .
COPY --from=builder /build/config ./config
COPY --from=builder /build/static ./static
COPY start.sh .
COPY debug-entrypoint.sh .

# Create non-root user and set up directories
RUN adduser -D -h /app muraena && \
    mkdir -p /app/certs && \
    chown -R muraena:muraena /app && \
    chmod +x start.sh && \
    chmod +x debug-entrypoint.sh

USER muraena

# Expose ports
EXPOSE 80 443

# Set environment variable for config file path with default value
ENV CONFIG_FILE=/app/config/local-ssl.toml

# Use the start script
ENTRYPOINT ["/app/start.sh"]
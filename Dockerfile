# Build stage
FROM golang:1.23-alpine AS builder

# Set necessary environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Install git and ca-certificates (needed for fetching modules)
RUN apk update && apk add --no-cache git ca-certificates tzdata

# Create a non-root user
RUN adduser -D -g '' appuser

# Set the working directory
WORKDIR /build

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

# Copy the source code
COPY . .

# Build the application
RUN go build -ldflags='-w -s -extldflags "-static"' -a -installsuffix cgo -o chat-agent ./cmd/server

# Final stage
FROM scratch

# Import from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd

# Copy the binary
COPY --from=builder /build/chat-agent /app/chat-agent

# Copy config files
COPY --from=builder /build/configs /app/configs

# Use non-root user
USER appuser

# Set working directory
WORKDIR /app

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/app/chat-agent", "health"] || exit 1

# Run the application
ENTRYPOINT ["/app/chat-agent"]


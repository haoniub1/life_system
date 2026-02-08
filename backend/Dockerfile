# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy source code
COPY . .

# Download dependencies and build
RUN go mod tidy && CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/life-system main.go

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Set timezone to Shanghai
ENV TZ=Asia/Shanghai

# Copy binary from builder
COPY --from=builder /app/life-system /app/life-system

# Create data directory
RUN mkdir -p /app/data /app/etc

# Copy config
COPY etc/config.example.yaml /app/etc/config.yaml

EXPOSE 8081

CMD ["/app/life-system", "-f", "/app/etc/config.yaml"]

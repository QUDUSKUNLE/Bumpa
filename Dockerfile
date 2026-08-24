# =========================
# Build stage
# =========================
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Required for fetching Go dependencies
RUN apk add --no-cache git ca-certificates

# Copy dependency files first for better Docker layer caching
COPY go.mod go.sum ./

RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/bumpa ./main.go


# =========================
# Runtime stage
# =========================
FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

# Copy compiled binary
COPY --from=builder /app/bumpa /app/bumpa

# Application port
EXPOSE 8081

# Start application
CMD ["/app/bumpa"]

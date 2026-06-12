# ── Build stage ───────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install git (needed by go mod download for some deps)
RUN apk add --no-cache git

# Download dependencies first (layer cache)
COPY go.mod go.sum ./
RUN go mod download

# Copy source and generate swagger docs
COPY . .
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    swag init -g cmd/main.go -o docs

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /cyberguard ./cmd/main.go

# ── Runtime stage ─────────────────────────────────────────────────────────────
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /cyberguard .

EXPOSE 8080

CMD ["/app/cyberguard"]

# ==============================================================
# Stage 1: Builder
# Compile the Go binary. This stage includes all dev tools.
# ==============================================================
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy dependency files and download modules (cached layer)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build a statically-linked binary with all optimizations
# -ldflags="-w -s" strips debug info to reduce binary size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server ./main.go

# ==============================================================
# Stage 2: Production Image
# Start from scratch to produce the smallest possible image (<50MB)
# ==============================================================
FROM scratch

# Import CA certificates and timezone data from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the compiled binary
COPY --from=builder /app/server /server

# Copy static assets and views (html templates)
COPY --from=builder /app/static /static
COPY --from=builder /app/views /views

# Copy environment file (if needed at runtime)
# COPY .env .env

EXPOSE 8081 8082

# Run the binary
ENTRYPOINT ["/server"]

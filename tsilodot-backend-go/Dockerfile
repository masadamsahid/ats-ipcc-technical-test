# Stage 1: Build
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

# Stage 2: Runtime
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs
# Copy .env.example as .env if .env doesn't exist (handled by docker-compose)
# COPY --from=builder /app/.env . 

# Expose port
EXPOSE 8080

# Command to run the application
CMD ["./main"]

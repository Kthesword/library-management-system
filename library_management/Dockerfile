# Multi-stage Dockerfile for Console-Based Library Management System

# Stage 1: Build the Go binary
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy module definition
COPY go.mod ./
RUN go mod download

# Copy application source code
COPY . .

# Build statically linked Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o library_management main.go

# Stage 2: Lightweight runtime image
FROM alpine:3.21

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/library_management /app/library_management

# Set execution entrypoint
ENTRYPOINT ["/app/library_management"]


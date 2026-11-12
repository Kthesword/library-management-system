# Root Dockerfile pointing to library_management source
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY library_management/go.mod ./
RUN go mod download

COPY library_management/ .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o library_management main.go

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/library_management /app/library_management

ENTRYPOINT ["/app/library_management"]


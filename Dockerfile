# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod download
RUN go build -o ssht main.go

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/ssht /app/ssht
COPY config.toml /app/config.toml

# Set entrypoint
ENTRYPOINT ["/app/ssht"]
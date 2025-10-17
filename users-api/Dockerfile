FROM golang:1.24.1-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Build a static Linux binary (disable CGO) for the server entrypoint
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the built binary from the builder stage. Do not copy .env by default;
# mount or provide environment variables at runtime instead.
COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]

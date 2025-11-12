# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS builder
WORKDIR /app

# Install git for modules that need it
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/api

FROM alpine:3.20
WORKDIR /app

ENV PORT=8081
EXPOSE 8081

COPY --from=builder /app/server /usr/local/bin/reservations-api

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 CMD wget -qO- http://127.0.0.1:${PORT}/health || exit 1

ENTRYPOINT ["/usr/local/bin/reservations-api"]

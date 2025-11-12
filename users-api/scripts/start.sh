#!/bin/bash

# Start script for users-api
# This script starts the users-api service in production mode

echo "🚀 Starting users-api..."

# Load environment variables if .env exists
if [ -f .env ]; then
    echo "📝 Loading environment variables from .env"
    export $(cat .env | grep -v '^#' | xargs)
fi

# Check if MySQL is running
echo "🔍 Checking MySQL connection..."
until nc -z ${DB_HOST:-localhost} ${DB_PORT:-3306}; do
    echo "⏳ Waiting for MySQL to be ready..."
    sleep 2
done
echo "✅ MySQL is ready"

# Build the application
echo "🔨 Building application..."
go build -o bin/users-api ./cmd/server/main.go

# Run the application
echo "▶️  Starting server on port 8080..."
./bin/users-api

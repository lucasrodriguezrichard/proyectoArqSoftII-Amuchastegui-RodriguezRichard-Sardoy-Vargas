#!/bin/bash

# Development script for users-api with hot-reload using Air
# This script starts the users-api service in development mode with automatic reloading

echo "🔧 Starting users-api in development mode (hot-reload)..."

# Load environment variables if .env exists
if [ -f .env ]; then
    echo "📝 Loading environment variables from .env"
    export $(cat .env | grep -v '^#' | xargs)
fi

# Check if Air is installed
if ! command -v air &> /dev/null; then
    echo "⚠️  Air is not installed. Installing..."
    go install github.com/cosmtrek/air@latest
fi

# Check if MySQL is running
echo "🔍 Checking MySQL connection..."
until nc -z ${DB_HOST:-localhost} ${DB_PORT:-3306}; do
    echo "⏳ Waiting for MySQL to be ready..."
    sleep 2
done
echo "✅ MySQL is ready"

# Create .air.toml if it doesn't exist
if [ ! -f .air.toml ]; then
    echo "📄 Creating .air.toml configuration..."
    cat > .air.toml << 'EOF'
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/server/main.go"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_error = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true
EOF
fi

# Run Air for hot-reload
echo "▶️  Starting server with hot-reload..."
air

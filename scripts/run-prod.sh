#!/bin/bash

set -e

# Always run from project root
cd "$(dirname "$0")/.."

echo "Ensuring logs/ exists..."
mkdir -p logs

echo "Copying .env file if missing..."
if [ ! -f "parser/.env" ] && [ -f "parser/.env.example" ]; then
  cp "parser/.env.example" "parser/.env"
  echo "Created parser/.env from example"
fi

echo "Starting production container..."
podman-compose -f podman-compose-prod.yml up --build -d

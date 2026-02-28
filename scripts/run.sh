#!/bin/bash

set -e

# Always run from project root
cd "$(dirname "$0")/.."

echo "Ensuring logs/ exists..."
mkdir -p logs

echo "Copying .env files if missing..."
for dir in parser webhook webhook-server; do
  if [ ! -f "$dir/.env" ] && [ -f "$dir/.env.example" ]; then
    cp "$dir/.env.example" "$dir/.env"
    echo "Created $dir/.env from example"
  fi
done

echo "Building services..."
make build

echo "Launching stack with podman-compose..."
podman-compose -f podman-compose.yml up --build

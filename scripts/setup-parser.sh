#!/bin/bash
set -e

# Always run from project root
cd "$(dirname "$0")/.."

echo "Verifying you are in the correct project root directory..."

REQUIRED_ITEMS=("parser" "podman-compose-prod.yml" "Makefile")

for item in "${REQUIRED_ITEMS[@]}"; do
  if [ ! -e "$item" ]; then
    echo "Missing required item: $item"
    echo "Please run this script from the root of the SMTPHook project directory."
    exit 1
  fi
done

# Detect and install system packages
echo "Detecting and installing required packages..."

if command -v apt-get &>/dev/null; then
  # Debian, Ubuntu, Linux Mint
  sudo apt update
  sudo apt install -y podman python3-pip pipx curl wget
elif command -v dnf &>/dev/null; then
  # Fedora
  sudo dnf install -y podman python3-pip pipx curl wget
elif command -v pacman &>/dev/null; then
  # Arch Linux
  sudo pacman -Sy --noconfirm podman python-pipx python-pip curl wget
elif command -v apk &>/dev/null; then
  # Alpine
  sudo apk add podman py3-pip pipx curl wget
else
  echo "Unsupported package manager or distribution."
  exit 1
fi

# Ensure pipx path
pipx ensurepath
export PATH="$HOME/.local/bin:$PATH"

# Install podman-compose if missing
echo "Ensuring podman-compose is installed..."
if ! command -v podman-compose &>/dev/null; then
  pipx install podman-compose
else
  echo "podman-compose is already installed."
fi

# Create .env from example if missing
if [ ! -f "parser/.env" ]; then
  if [ -f "parser/.env.production.example" ]; then
    cp parser/.env.production.example parser/.env
    echo "Created parser/.env from .env.production.example"
  else
    echo "Missing parser/.env.production.example"
    exit 1
  fi
else
  echo "parser/.env already exists."
fi

# Prompt user to review .env
echo ""
echo "Please review and edit the parser/.env file before running the container."
echo "It should look something like:"
echo ""
echo "  POLL_INTERVAL=5"
echo "  WEBHOOK_URL=https://your-api.example.com/webhook"
echo "  MAIL_DIR=/mail/inbox"
echo ""
read -p "Do you want to open parser/.env in nano now? [Y/n]: " EDIT_ENV

if [[ "$EDIT_ENV" =~ ^[Yy]$ || -z "$EDIT_ENV" ]]; then
  nano parser/.env
fi

# Create folders if needed
echo "Creating default mail inbox and logs directory if missing..."
mkdir -p mail/inbox logs

echo ""
echo "Container-only setup complete."
echo ""
echo "To start the parser container:"
echo "  podman-compose -f podman-compose-prod.yml up -d"

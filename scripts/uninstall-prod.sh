#!/bin/bash
set -e

# Always run from project root
cd "$(dirname "$0")/.."

echo "🛑 Stopping and disabling production services..."

systemctl --user stop container-parser-prod.container || true
systemctl --user disable container-parser-prod.container || true

echo "🧹 Cleaning Quadlet config..."
rm -f ~/.config/containers/systemd/container-parser-prod.container

systemctl --user daemon-reload
echo "✅ Uninstalled production systemd container unit"

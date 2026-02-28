#!/bin/bash
set -e

# Always run from project root
cd "$(dirname "$0")/.."

echo "🛑 Stopping services..."
sudo systemctl stop smtphook.target || true
sudo systemctl stop parser.service || true
sudo systemctl stop webhook.service || true
sudo systemctl stop webhook-server.service || true

echo "Disabling services..."
sudo systemctl disable smtphook.target || true
sudo systemctl disable parser.service || true
sudo systemctl disable webhook.service || true
sudo systemctl disable webhook-server.service || true

echo "Removing systemd unit files..."
sudo rm -f /etc/systemd/system/parser.service
sudo rm -f /etc/systemd/system/webhook.service
sudo rm -f /etc/systemd/system/webhook-server.service
sudo rm -f /etc/systemd/system/smtphook.target
sudo systemctl daemon-reload

echo "Removing installed binaries and service folders..."
sudo rm -rf /opt/smtphook/bin
sudo rm -rf /opt/smtphook/parser
sudo rm -rf /opt/smtphook/webhook
sudo rm -rf /opt/smtphook/webhook-server

echo "Removing logrotate config..."
sudo rm -f /etc/logrotate.d/smtphook

echo "SMTPHook has been uninstalled."

read -p "Do you want to delete the local logs/ folder as well? [y/N] " choice
case "$choice" in 
  y|Y ) rm -rf logs && echo "logs/ deleted.";;
  * ) echo "logs/ kept.";;
esac

#!/bin/bash
set -e

# Always run from project root
cd "$(dirname "$0")/.."

echo "🔎 Running SMTPHook production diagnostic..."
echo ""

# Check parser binary only
echo "Checking parser binary..."
if [ -x "/opt/smtphook/bin/parser" ]; then
  echo "/opt/smtphook/bin/parser exists"
else
  echo "/opt/smtphook/bin/parser missing"
fi

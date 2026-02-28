#!/bin/bash
set -e

# Always run from project root
SCRIPT_DIR="$(dirname "$0")"
cd "$SCRIPT_DIR/.."

COMMAND="$1"

if [ -z "$COMMAND" ]; then
  echo "Usage: $0 {setup|run|reset|uninstall|diagnose}"
  exit 1
fi

case "$COMMAND" in
  setup)
    "$SCRIPT_DIR/setup-production.sh"
    ;;
  run)
    "$SCRIPT_DIR/run-prod.sh"
    ;;
  reset)
    "$SCRIPT_DIR/reset-prod.sh"
    ;;
  uninstall)
    "$SCRIPT_DIR/uninstall-prod.sh"
    ;;
  diagnose)
    "$SCRIPT_DIR/diagnose-prod.sh"
    ;;
  *)
    echo "❌ Unknown command: $COMMAND"
    echo "Usage: $0 {setup|run|reset|uninstall|diagnose}"
    exit 1
    ;;
esac

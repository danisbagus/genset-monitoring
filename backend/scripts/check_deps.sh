#!/usr/bin/env bash
# ──────────────────────────────────────────────────────────
# scripts/check_deps.sh
# Validates that required CLI tools are installed.
# ──────────────────────────────────────────────────────────
set -euo pipefail

required=(go docker docker-compose swag)

for cmd in "${required[@]}"; do
  if ! command -v "$cmd" &>/dev/null; then
    echo "❌  Missing required tool: $cmd"
    exit 1
  fi
  echo "✅  $cmd found: $(command -v "$cmd")"
done

echo ""
echo "All required tools are present."

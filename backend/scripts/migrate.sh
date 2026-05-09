#!/usr/bin/env bash
# ══════════════════════════════════════════════════════════════════
# scripts/migrate.sh
# Thin wrapper around golang-migrate for database migrations.
#
# Usage:
#   ./scripts/migrate.sh up
#   ./scripts/migrate.sh down
#   ./scripts/migrate.sh version
#   ./scripts/migrate.sh force <version>
#   ./scripts/migrate.sh create <name>
#
# Reads DATABASE_URL from .env if not already set in environment.
# ══════════════════════════════════════════════════════════════════
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$(dirname "$SCRIPT_DIR")"
MIGRATIONS_DIR="${BACKEND_DIR}/migrations"
ENV_FILE="${BACKEND_DIR}/.env"

# ── Load .env if DATABASE_URL is not already in the environment ───
if [[ -z "${DATABASE_URL:-}" ]] && [[ -f "$ENV_FILE" ]]; then
  # Parse key=value lines, ignoring comments and blanks
  while IFS='=' read -r key value; do
    [[ "$key" =~ ^#.*$ || -z "$key" ]] && continue
    key="${key%%#*}"          # strip inline comments
    key="${key%"${key##*[![:space:]]}"}"  # trim trailing spaces
    value="${value%%#*}"
    value="${value%"${value##*[![:space:]]}"}"
    export "$key"="$value"
  done < "$ENV_FILE"
fi

# ── Build DATABASE_URL from individual DB_* vars if not set ───────
if [[ -z "${DATABASE_URL:-}" ]]; then
  DB_HOST="${DB_HOST:-localhost}"
  DB_PORT="${DB_PORT:-5432}"
  DB_USER="${DB_USER:-postgres}"
  DB_PASSWORD="${DB_PASSWORD:-postgres}"
  DB_NAME="${DB_NAME:-genset_monitoring}"
  DB_SSL_MODE="${DB_SSL_MODE:-disable}"
  DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}"
fi

# ── Check that migrate binary is available ────────────────────────
if ! command -v migrate &>/dev/null; then
  echo "❌  'migrate' CLI not found."
  echo "    Install with:"
  echo "      go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
  echo "    Or via the Makefile: make migrate-install"
  exit 1
fi

CMD="${1:-help}"

case "$CMD" in
  up)
    echo "▶  Running migrations UP..."
    migrate -database "$DATABASE_URL" -path "$MIGRATIONS_DIR" up
    echo "✅  Migrations applied successfully."
    ;;
  down)
    STEPS="${2:-1}"
    echo "▼  Rolling back ${STEPS} migration(s)..."
    migrate -database "$DATABASE_URL" -path "$MIGRATIONS_DIR" down "$STEPS"
    echo "✅  Rollback complete."
    ;;
  version)
    echo "ℹ️   Current migration version:"
    migrate -database "$DATABASE_URL" -path "$MIGRATIONS_DIR" version
    ;;
  force)
    VERSION="${2:?Usage: $0 force <version>}"
    echo "⚠️   Forcing migration version to ${VERSION}..."
    migrate -database "$DATABASE_URL" -path "$MIGRATIONS_DIR" force "$VERSION"
    echo "✅  Version forced to ${VERSION}."
    ;;
  create)
    NAME="${2:?Usage: $0 create <migration_name>}"
    echo "📝  Creating new migration: ${NAME}..."
    migrate create -ext sql -dir "$MIGRATIONS_DIR" -seq "$NAME"
    echo "✅  Migration files created in ${MIGRATIONS_DIR}."
    ;;
  drop)
    echo "🔥  Dropping ALL tables (irreversible)..."
    migrate -database "$DATABASE_URL" -path "$MIGRATIONS_DIR" drop -f
    echo "✅  All tables dropped."
    ;;
  help|*)
    echo ""
    echo "  Usage: $0 <command> [args]"
    echo ""
    echo "  Commands:"
    echo "    up              Apply all pending migrations"
    echo "    down [N]        Roll back N migrations (default: 1)"
    echo "    version         Show current migration version"
    echo "    force <version> Force-set migration version (use after manual fix)"
    echo "    create <name>   Create a new up/down migration pair"
    echo "    drop            Drop all tables (DESTRUCTIVE)"
    echo ""
    ;;
esac

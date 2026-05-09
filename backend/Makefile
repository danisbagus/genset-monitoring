# ══════════════════════════════════════════════════════════════════════
#  Genset Monitoring — Makefile
# ══════════════════════════════════════════════════════════════════════

.PHONY: all run dev build swagger tidy test lint clean \
        migrate-install migrate-up migrate-down migrate-version \
        migrate-force migrate-create migrate-drop \
        docker-up docker-down docker-logs help

BINARY      := server
CMD_DIR     := ./cmd/api
OUTPUT_DIR  := ./bin
MAIN        := $(CMD_DIR)/main.go
ENV_FILE    := .env
MIGRATE_SH  := ./scripts/migrate.sh

# ── Default ────────────────────────────────────────────────────────
all: build

# ════════════════════════════════════════════════════════════════
#  Application
# ════════════════════════════════════════════════════════════════

## run: Build and run the production binary
run: build
	ENV_FILE=$(ENV_FILE) $(OUTPUT_DIR)/$(BINARY)

## dev: Run with hot-reload (air) — falls back to go run
dev:
	@if command -v air > /dev/null 2>&1; then \
		air; \
	else \
		echo "⚠️  air not found, using go run (install air for hot-reload)"; \
		ENV_FILE=$(ENV_FILE) go run $(MAIN); \
	fi

## build: Compile the binary
build:
	@mkdir -p $(OUTPUT_DIR)
	go build -ldflags="-w -s" -o $(OUTPUT_DIR)/$(BINARY) $(CMD_DIR)
	@echo "✅  Binary built: $(OUTPUT_DIR)/$(BINARY)"

# ════════════════════════════════════════════════════════════════
#  Database Migrations  (golang-migrate)
# ════════════════════════════════════════════════════════════════

## migrate-install: Install the golang-migrate CLI
migrate-install:
	@echo "Installing golang-migrate CLI..."
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "✅  migrate CLI installed"

## migrate-up: Apply all pending migrations
migrate-up:
	@bash $(MIGRATE_SH) up

## migrate-down: Roll back 1 migration  (use N=n for multiple: make migrate-down N=3)
migrate-down:
	@bash $(MIGRATE_SH) down $(N)

## migrate-version: Show the current migration version
migrate-version:
	@bash $(MIGRATE_SH) version

## migrate-force: Force migration version  (use V=<version>: make migrate-force V=1)
migrate-force:
	@bash $(MIGRATE_SH) force $(V)

## migrate-create: Create a new migration  (use NAME=<name>: make migrate-create NAME=add_sessions)
migrate-create:
	@bash $(MIGRATE_SH) create $(NAME)

## migrate-drop: Drop ALL tables — DESTRUCTIVE, use with caution
migrate-drop:
	@bash $(MIGRATE_SH) drop

# ════════════════════════════════════════════════════════════════
#  Code Quality
# ════════════════════════════════════════════════════════════════

## swagger: Generate Swagger/OpenAPI docs via swag
swagger:
	@if ! command -v swag > /dev/null 2>&1; then \
		echo "Installing swag..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
	fi
	swag init \
		--generalInfo cmd/api/main.go \
		--output docs \
		--parseDependency \
		--parseInternal
	@echo "✅  Swagger docs generated in ./docs"

## tidy: Tidy go.mod / go.sum
tidy:
	go mod tidy
	@echo "✅  go.mod and go.sum tidied"

## test: Run tests with race detector and coverage report
test:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅  Tests done. See coverage.html"

## lint: Run golangci-lint
lint:
	@if ! command -v golangci-lint > /dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	golangci-lint run ./...

# ════════════════════════════════════════════════════════════════
#  Docker
# ════════════════════════════════════════════════════════════════

## docker-up: Start all infrastructure and API containers
docker-up:
	docker-compose -f ../docker-compose.yml up -d
	@echo "✅  Infrastructure and API containers started"

## docker-down: Stop all containers
docker-down:
	docker-compose -f ../docker-compose.yml down
	@echo "✅  Infrastructure and API containers stopped"

## docker-logs: Tail logs from all containers
docker-logs:
	docker-compose -f ../docker-compose.yml logs -f

# ════════════════════════════════════════════════════════════════
#  Cleanup
# ════════════════════════════════════════════════════════════════

## clean: Remove build artifacts and coverage reports
clean:
	rm -rf $(OUTPUT_DIR) coverage.out coverage.html
	@echo "✅  Build artifacts cleaned"

# ════════════════════════════════════════════════════════════════
#  Help
# ════════════════════════════════════════════════════════════════

## help: List all available make commands
help:
	@echo ""
	@echo "  Genset Monitoring — Available Commands"
	@echo "  ═══════════════════════════════════════"
	@grep -E '^## ' Makefile | sed 's/## /  /' | column -t -s ':'
	@echo ""

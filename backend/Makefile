# ══════════════════════════════════════════════════════════════
#  Genset Monitoring — Makefile
# ══════════════════════════════════════════════════════════════

.PHONY: all run dev build swagger tidy test lint clean docker-up docker-down help

BINARY      := server
CMD_DIR     := ./cmd/api
OUTPUT_DIR  := ./bin
MAIN        := $(CMD_DIR)/main.go
ENV_FILE    := .env

# ── Default ────────────────────────────────────────────────
all: build

# ── Run (production binary) ────────────────────────────────
run: build
	ENV_FILE=$(ENV_FILE) $(OUTPUT_DIR)/$(BINARY)

# ── Dev (hot-reload with air, falls back to go run) ────────
dev:
	@if command -v air > /dev/null 2>&1; then \
		air; \
	else \
		echo "⚠️  air not found, using go run (install air for hot-reload)"; \
		ENV_FILE=$(ENV_FILE) go run $(MAIN); \
	fi

# ── Build binary ───────────────────────────────────────────
build:
	@mkdir -p $(OUTPUT_DIR)
	go build -ldflags="-w -s" -o $(OUTPUT_DIR)/$(BINARY) $(CMD_DIR)
	@echo "✅  Binary built: $(OUTPUT_DIR)/$(BINARY)"

# ── Generate Swagger docs ──────────────────────────────────
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

# ── Tidy dependencies ──────────────────────────────────────
tidy:
	go mod tidy
	@echo "✅  go.mod and go.sum tidied"

# ── Run tests ──────────────────────────────────────────────
test:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅  Tests done. See coverage.html"

# ── Lint (requires golangci-lint) ──────────────────────────
lint:
	@if ! command -v golangci-lint > /dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	golangci-lint run ./...

# ── Docker ─────────────────────────────────────────────────
docker-up:
	docker-compose -f ../docker-compose.yml up -d
	@echo "✅  Infrastructure and API containers started"

docker-down:
	docker-compose -f ../docker-compose.yml down
	@echo "✅  Infrastructure and API containers stopped"

docker-logs:
	docker-compose -f ../docker-compose.yml logs -f

# ── Clean ──────────────────────────────────────────────────
clean:
	rm -rf $(OUTPUT_DIR) coverage.out coverage.html
	@echo "✅  Build artifacts cleaned"

# ── Help ───────────────────────────────────────────────────
help:
	@echo ""
	@echo "  Genset Monitoring — Available Commands"
	@echo "  ─────────────────────────────────────"
	@echo "  make run          Build and run the server"
	@echo "  make dev          Run with hot-reload (air) or go run"
	@echo "  make build        Compile the binary"
	@echo "  make swagger      Generate Swagger docs via swag init"
	@echo "  make tidy         Tidy go.mod / go.sum"
	@echo "  make test         Run tests with coverage"
	@echo "  make lint         Run golangci-lint"
	@echo "  make docker-up    Start infrastructure (postgres/redis/mqtt)"
	@echo "  make docker-down  Stop infrastructure containers"
	@echo "  make clean        Remove build artifacts"
	@echo ""

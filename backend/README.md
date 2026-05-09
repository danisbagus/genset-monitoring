# Genset Monitoring — Backend

Production-ready IoT Genset Monitoring backend built with Go.

## Tech Stack

| Layer       | Technology                     |
|-------------|-------------------------------|
| HTTP        | Gin                           |
| ORM         | GORM + pgx                   |
| Database    | PostgreSQL 16                 |
| Cache       | Redis 7                       |
| Messaging   | MQTT (Paho) + Mosquitto       |
| WebSocket   | Gorilla WebSocket             |
| Auth        | JWT (golang-jwt/jwt v5)       |
| Logging     | Uber Zap                      |
| Config      | Viper + godotenv              |
| Docs        | Swaggo / Swagger UI           |

---

## Quick Start

### Prerequisites

- Go 1.22+
- Docker & Docker Compose
- `swag` CLI (`go install github.com/swaggo/swag/cmd/swag@latest`)

### 1. Clone & configure

```bash
cp .env.example .env
# Edit .env with your values
```

### 2. Start infrastructure

```bash
make docker-up
```

### 3. Run the server

```bash
make dev
```

Server starts at **http://localhost:8080**

---

## API Documentation

Swagger UI is available at:

```
http://localhost:8080/swagger/index.html
```

Regenerate docs after changing annotations:

```bash
make swagger
```

---

## Healthcheck

```
GET /api/v1/health
```

```json
{
  "success": true,
  "message": "service healthy",
  "data": {
    "postgres": "connected",
    "redis": "connected",
    "mqtt": "connected"
  }
}
```

---

## Makefile Commands

| Command           | Description                            |
|-------------------|----------------------------------------|
| `make run`        | Build and run the binary               |
| `make dev`        | Hot-reload (air) or go run             |
| `make build`      | Compile binary to `./bin/server`       |
| `make swagger`    | Generate Swagger docs                  |
| `make tidy`       | Tidy go.mod / go.sum                   |
| `make test`       | Run tests with coverage report         |
| `make lint`       | Run golangci-lint                      |
| `make docker-up`  | Start postgres/redis/mqtt containers   |
| `make docker-down`| Stop infrastructure containers         |

---

## Project Structure

```
backend/
├── cmd/api/main.go              # Entry point, DI wiring, graceful shutdown
├── docs/                        # Swagger generated docs (swag init)
├── internal/
│   ├── config/                  # Viper config loader
│   ├── handler/                 # HTTP handlers (Gin)
│   ├── service/                 # Business logic interfaces + implementations
│   ├── repository/              # Database access layer (GORM)
│   ├── middleware/              # Gin middlewares (auth, logger, recovery)
│   ├── websocket/               # Gorilla WebSocket hub
│   ├── mqtt/                    # MQTT service interface
│   ├── auth/                    # Auth service interface (placeholder)
│   ├── model/                   # GORM domain models
│   └── infrastructure/
│       ├── database/            # PostgreSQL singleton
│       ├── redis/               # Redis singleton
│       ├── mqtt/                # MQTT singleton with auto-reconnect
│       └── logger/              # Zap logger setup
├── pkg/
│   ├── response/                # Reusable HTTP response helpers
│   ├── utils/                   # UUID, time, crypto utilities
│   └── validator/               # Struct validation wrapper
├── docker/
│   └── mosquitto/config/        # Mosquitto MQTT config
├── scripts/                     # Helper shell scripts
├── .env.example                 # Environment variable template
├── docker-compose.yml           # Infrastructure stack
├── Dockerfile                   # Multi-stage build
└── Makefile                     # Developer commands
```

---

## Architecture

The project follows **Clean Architecture** principles:

- **Handler** — HTTP request/response, no business logic
- **Service** — Business logic, implements interfaces
- **Repository** — Data access, wraps GORM operations
- **Infrastructure** — External connections (DB, Redis, MQTT)
- **Model** — Domain entities (GORM models)

Dependencies flow inward: `handler → service → repository → infrastructure`

---

## WebSocket

Connect to the WebSocket endpoint:

```
ws://localhost:8080/ws?client_id=dashboard-1
```

The server broadcasts telemetry updates to all connected clients via the Hub.

---

## Environment Variables

See [`.env.example`](.env.example) for a full reference.

Key variables:

| Variable          | Default               | Description              |
|-------------------|-----------------------|--------------------------|
| `APP_PORT`        | `8080`                | HTTP server port         |
| `APP_ENV`         | `development`         | `development|production` |
| `DB_HOST`         | `localhost`           | PostgreSQL host          |
| `REDIS_HOST`      | `localhost`           | Redis host               |
| `MQTT_BROKER`     | `localhost`           | MQTT broker host         |
| `JWT_SECRET`      | *(change me)*         | JWT signing secret       |

---

## License

Apache 2.0

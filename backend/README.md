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
| Auth        | JWT + Refresh Token Rotation  |
| Security    | Bcrypt + SHA-256 Hashing      |
| Logging     | Uber Zap                      |
| Config      | Viper + godotenv              |
| Docs        | Swaggo / Swagger UI           |

---

## Quick Start

### Prerequisites

- Go 1.22+
- Docker & Docker Compose
- `golang-migrate` CLI (`make migrate-install`)
- `swag` CLI (`go install github.com/swaggo/swag/cmd/swag@latest`)

### 1. Clone & configure

```bash
cp .env.example .env
# Edit .env with your values, especially JWT_SECRET
```

### 2. Start infrastructure

```bash
make docker-up
```

### 3. Run database migrations

```bash
make migrate-up
```

### 4. Run the server

```bash
make dev
```

Server starts at **http://localhost:8080**

---

## Authentication Flow

This project implements a secure **JWT + Refresh Token Rotation** strategy:

1.  **Login**: User provides credentials. Server returns `access_token` (short-lived) and `refresh_token` (long-lived).
2.  **Access**: Client uses `access_token` in `Authorization: Bearer <token>` header.
3.  **Refresh**: When `access_token` expires, client sends `refresh_token` to `/auth/refresh`.
4.  **Rotation**: Server revokes the old `refresh_token`, issues a **new pair** of tokens. This prevents replay attacks.
5.  **Security**: Refresh tokens are stored as **SHA-256 hashes** in the database. Raw tokens are never persisted. Passwords use **Bcrypt** (cost 12).

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

## Database Migrations

We use `golang-migrate` for versioned database schema management.

| Command                  | Description                                      |
|--------------------------|--------------------------------------------------|
| `make migrate-up`        | Apply all pending migrations                     |
| `make migrate-down`      | Roll back the last migration                     |
| `make migrate-create NAME=x` | Create a new migration file pair              |
| `make migrate-version`   | Show current schema version                      |
| `make migrate-install`   | Install the `migrate` CLI tool                   |

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
├── cmd/api/main.go              # Entry point, DI wiring, routes
├── migrations/                  # SQL migration files (up/down)
├── internal/
│   ├── handler/                 # HTTP handlers (Auth, Health)
│   ├── service/                 # Business logic (AuthService)
│   ├── repository/              # Data access (UserRepository, TokenRepository)
│   ├── model/                   # GORM models (User, RefreshToken)
│   ├── middleware/              # Auth, Logger, Recovery middlewares
│   ├── infrastructure/          # DB, Redis, MQTT, Logger drivers
│   └── config/                  # Configuration loader
├── pkg/
│   ├── jwtutil/                 # JWT manager (sign/parse)
│   ├── hashutil/                # Password (bcrypt) & Token (sha256) utils
│   ├── response/                # Standard JSON response helpers
│   └── validator/               # Request validation wrapper
└── scripts/                     # Migration wrapper scripts
```

---

## Environment Variables

Key variables in `.env`:

| Variable                 | Default       | Description                       |
|--------------------------|---------------|-----------------------------------|
| `JWT_SECRET`             | *(required)*  | Secret key for signing tokens     |
| `JWT_EXPIRATION`         | `1h`          | Access token TTL                  |
| `JWT_REFRESH_EXPIRATION` | `168h`        | Refresh token TTL (7 days)        |
| `DB_PORT`                | `5435`        | PostgreSQL port (docker-compose)  |

---

## License

Apache 2.0

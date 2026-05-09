// Package auth is reserved for future domain-specific auth utilities.
// The full authentication implementation lives in:
//   - internal/service/auth.go      — AuthService interface + implementation
//   - internal/repository/user.go   — UserRepository
//   - internal/repository/refresh_token.go — RefreshTokenRepository
//   - internal/handler/auth.go      — HTTP handlers
//   - internal/middleware/auth.go   — JWT middleware + RBAC
//   - pkg/jwtutil/jwt.go            — JWT token manager
//   - pkg/hashutil/hash.go          — bcrypt + SHA-256 utilities
package auth

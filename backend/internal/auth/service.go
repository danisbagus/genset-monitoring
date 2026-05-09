package auth

// AuthService defines the contract for authentication operations.
// Full implementation is outside the MVP scope.
type AuthService interface {
	// Login authenticates a user and returns a signed JWT token.
	Login(ctx interface{}, email, password string) (string, error)

	// Register creates a new user account.
	Register(ctx interface{}, name, email, password string) error
}

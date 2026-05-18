package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/danisbagus/genset-monitoring/backend/internal/model"
	"github.com/danisbagus/genset-monitoring/backend/internal/repository"
	"github.com/danisbagus/genset-monitoring/backend/pkg/hashutil"
	"github.com/danisbagus/genset-monitoring/backend/pkg/jwtutil"
)

// ── DTOs ─────────────────────────────────────────────────────────

// RegisterInput is the payload required to create a new user account.
type RegisterInput struct {
	Username string
	Email    string
	Password string
	Role     model.UserRole // defaults to RoleViewer if empty
}

// RegisterOutput is returned after a successful registration.
type RegisterOutput struct {
	ID       uuid.UUID      `json:"id"`
	Username string         `json:"username"`
	Email    string         `json:"email"`
	Role     model.UserRole `json:"role"`
}

// LoginInput is the payload required for authentication.
type LoginInput struct {
	Username string
	Password string
}

// LoginOutput is returned after a successful login.
type LoginOutput struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    int64       `json:"expires_in"`
	User         UserProfile `json:"user"`
}

// RefreshOutput is returned after a successful token refresh.
type RefreshOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// UserProfile is a public-facing user representation (no sensitive fields).
type UserProfile struct {
	ID       uuid.UUID      `json:"id"`
	Username string         `json:"username"`
	Email    string         `json:"email"`
	Role     model.UserRole `json:"role"`
}

// ── Interface ─────────────────────────────────────────────────────

// AuthService defines the contract for authentication operations.
type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (*RegisterOutput, error)
	Login(ctx context.Context, input LoginInput, userAgent, ipAddress string) (*LoginOutput, error)
	RefreshTokens(ctx context.Context, rawRefreshToken, userAgent, ipAddress string) (*RefreshOutput, error)
	Me(ctx context.Context, userID uuid.UUID) (*UserProfile, error)
	Logout(ctx context.Context, rawRefreshToken string) error
}

// ── Implementation ────────────────────────────────────────────────

type authService struct {
	userRepo   repository.UserRepository
	tokenRepo  repository.RefreshTokenRepository
	jwtManager *jwtutil.Manager
	refreshTTL time.Duration
	log        *zap.Logger
}

// NewAuthService constructs an AuthService with its dependencies.
func NewAuthService(
	userRepo repository.UserRepository,
	tokenRepo repository.RefreshTokenRepository,
	jwtManager *jwtutil.Manager,
	refreshTTL time.Duration,
	log *zap.Logger,
) AuthService {
	return &authService{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		jwtManager: jwtManager,
		refreshTTL: refreshTTL,
		log:        log,
	}
}

// Register creates a new user account.
// Returns ErrUsernameExists or ErrEmailExists on duplicates.
func (s *authService) Register(ctx context.Context, input RegisterInput) (*RegisterOutput, error) {
	// Validate unique username
	exists, err := s.userRepo.ExistsByUsername(ctx, input.Username)
	if err != nil {
		return nil, fmt.Errorf("authService.Register: %w", err)
	}
	if exists {
		return nil, ErrUsernameExists
	}

	// Validate unique email
	exists, err = s.userRepo.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return nil, fmt.Errorf("authService.Register: %w", err)
	}
	if exists {
		return nil, ErrEmailExists
	}

	// Hash password
	hash, err := hashutil.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("authService.Register: %w", err)
	}

	role := input.Role
	if role == "" {
		role = model.RoleViewer
	}

	user := &model.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hash,
		Role:         role,
		IsActive:     true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("authService.Register: %w", err)
	}

	s.log.Info("user registered",
		zap.String("user_id", user.ID.String()),
		zap.String("username", user.Username),
		zap.String("role", string(user.Role)),
	)

	return &RegisterOutput{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}

// Login authenticates a user and issues a JWT access token + refresh token pair.
func (s *authService) Login(ctx context.Context, input LoginInput, userAgent, ipAddress string) (*LoginOutput, error) {
	user, err := s.userRepo.FindByUsername(ctx, input.Username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("authService.Login: %w", err)
	}

	if !user.IsActive {
		return nil, ErrAccountInactive
	}

	if err := hashutil.CheckPassword(input.Password, user.PasswordHash); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Issue access token
	accessToken, err := s.jwtManager.Generate(user.ID.String(), user.Username, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("authService.Login: %w", err)
	}

	// Issue refresh token
	rawToken, tokenHash, err := hashutil.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("authService.Login: %w", err)
	}

	rt := &model.RefreshToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().UTC().Add(s.refreshTTL),
		UserAgent: userAgent,
		IPAddress: ipAddress,
	}
	if err := s.tokenRepo.Create(ctx, rt); err != nil {
		return nil, fmt.Errorf("authService.Login: %w", err)
	}

	s.log.Info("user logged in",
		zap.String("user_id", user.ID.String()),
		zap.String("username", user.Username),
		zap.String("ip", ipAddress),
	)

	return &LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: rawToken,
		ExpiresIn:    s.jwtManager.ExpirationSeconds(),
		User: UserProfile{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		},
	}, nil
}

// RefreshTokens validates an existing refresh token, revokes it (rotation),
// and issues a brand-new access + refresh token pair.
func (s *authService) RefreshTokens(ctx context.Context, rawRefreshToken, userAgent, ipAddress string) (*RefreshOutput, error) {
	tokenHash := hashutil.HashSHA256(rawRefreshToken)

	rt, err := s.tokenRepo.FindByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, repository.ErrRefreshTokenNotFound) {
			return nil, ErrRefreshTokenInvalid
		}
		return nil, fmt.Errorf("authService.RefreshTokens: %w", err)
	}

	if !rt.IsValid() {
		return nil, ErrRefreshTokenInvalid
	}

	// Fetch the associated user
	user, err := s.userRepo.FindByID(ctx, rt.UserID)
	if err != nil {
		return nil, fmt.Errorf("authService.RefreshTokens: %w", err)
	}

	if !user.IsActive {
		return nil, ErrAccountInactive
	}

	// ── Token Rotation ────────────────────────────────────────────
	// Revoke the old refresh token BEFORE issuing the new one.
	// If the old token was already revoked (replay attack), we abort.
	if err := s.tokenRepo.RevokeByHash(ctx, tokenHash); err != nil {
		return nil, fmt.Errorf("authService.RefreshTokens: failed to revoke old token: %w", err)
	}

	// Issue new access token
	newAccessToken, err := s.jwtManager.Generate(user.ID.String(), user.Username, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("authService.RefreshTokens: %w", err)
	}

	// Issue new refresh token
	newRawToken, newTokenHash, err := hashutil.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("authService.RefreshTokens: %w", err)
	}

	newRT := &model.RefreshToken{
		UserID:    user.ID,
		TokenHash: newTokenHash,
		ExpiresAt: time.Now().UTC().Add(s.refreshTTL),
		UserAgent: userAgent,
		IPAddress: ipAddress,
	}
	if err := s.tokenRepo.Create(ctx, newRT); err != nil {
		return nil, fmt.Errorf("authService.RefreshTokens: %w", err)
	}

	return &RefreshOutput{
		AccessToken:  newAccessToken,
		RefreshToken: newRawToken,
		ExpiresIn:    s.jwtManager.ExpirationSeconds(),
	}, nil
}

// Me returns the public profile of an authenticated user.
func (s *authService) Me(ctx context.Context, userID uuid.UUID) (*UserProfile, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("authService.Me: %w", err)
	}

	return &UserProfile{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}

// Logout revokes a single refresh token (single-device logout).
// Use tokenRepo.RevokeAllByUserID for logout-all-devices.
func (s *authService) Logout(ctx context.Context, rawRefreshToken string) error {
	tokenHash := hashutil.HashSHA256(rawRefreshToken)
	if err := s.tokenRepo.RevokeByHash(ctx, tokenHash); err != nil {
		return fmt.Errorf("authService.Logout: %w", err)
	}
	return nil
}

// ── Sentinel errors ───────────────────────────────────────────────

var (
	ErrInvalidCredentials  = errors.New("invalid username or password")
	ErrUsernameExists      = errors.New("username already taken")
	ErrEmailExists         = errors.New("email already registered")
	ErrAccountInactive     = errors.New("account is inactive")
	ErrRefreshTokenInvalid = errors.New("refresh token is invalid or expired")
	ErrUserNotFound        = errors.New("user not found")
)

package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/danisbagus/genset-monitoring/backend/internal/middleware"
	"github.com/danisbagus/genset-monitoring/backend/internal/model"
	"github.com/danisbagus/genset-monitoring/backend/internal/service"
	"github.com/danisbagus/genset-monitoring/backend/pkg/response"
	"github.com/danisbagus/genset-monitoring/backend/pkg/validator"
)

// ── Request DTOs ─────────────────────────────────────────────────

// RegisterRequest is the expected JSON body for POST /api/v1/auth/register.
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=100"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=128"`
	Role     string `json:"role"     validate:"omitempty,oneof=admin operator viewer"`
}

// LoginRequest is the expected JSON body for POST /api/v1/auth/login.
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// RefreshRequest is the expected JSON body for POST /api/v1/auth/refresh and logout.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ── Handler ───────────────────────────────────────────────────────

// AuthHandler handles all authentication-related HTTP endpoints.
type AuthHandler struct {
	authSvc   service.AuthService
	validator *validator.Validator
	log       *zap.Logger
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authSvc service.AuthService, v *validator.Validator, log *zap.Logger) *AuthHandler {
	return &AuthHandler{
		authSvc:   authSvc,
		validator: v,
		log:       log,
	}
}

// Register creates a new user account.
//
//	@Summary		Register
//	@Description	Creates a new user account with the given credentials
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		RegisterRequest									true	"Registration payload"
//	@Success		201		{object}	response.Response{data=service.RegisterOutput}	"User created"
//	@Failure		400		{object}	response.ErrorResponse							"Validation error"
//	@Failure		409		{object}	response.ErrorResponse							"Username or email already taken"
//	@Failure		500		{object}	response.ErrorResponse							"Internal server error"
//	@Router			/api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", nil)
		return
	}

	if errs := h.validator.Validate(req); errs != nil {
		response.UnprocessableEntity(c, "validation failed", errs)
		return
	}

	role := model.RoleViewer
	if req.Role != "" {
		role = model.UserRole(req.Role)
	}

	out, err := h.authSvc.Register(c.Request.Context(), service.RegisterInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     role,
	})
	if err != nil {
		h.handleAuthError(c, err)
		return
	}

	response.Created(c, "user registered successfully", out)
}

// Login authenticates a user and returns JWT access + refresh tokens.
//
//	@Summary		Login
//	@Description	Authenticates a user and returns access and refresh tokens
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		LoginRequest								true	"Login payload"
//	@Success		200		{object}	response.Response{data=service.LoginOutput}	"Login successful"
//	@Failure		400		{object}	response.ErrorResponse						"Validation error"
//	@Failure		401		{object}	response.ErrorResponse						"Invalid credentials"
//	@Failure		500		{object}	response.ErrorResponse						"Internal server error"
//	@Router			/api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", nil)
		return
	}

	if errs := h.validator.Validate(req); errs != nil {
		response.UnprocessableEntity(c, "validation failed", errs)
		return
	}

	out, err := h.authSvc.Login(
		c.Request.Context(),
		service.LoginInput{Username: req.Username, Password: req.Password},
		c.GetHeader("User-Agent"),
		c.ClientIP(),
	)
	if err != nil {
		h.handleAuthError(c, err)
		return
	}

	response.OK(c, "login successful", out)
}

// RefreshToken rotates a refresh token and issues a new token pair.
//
//	@Summary		Refresh Token
//	@Description	Rotates a refresh token and returns a new access + refresh token pair
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		RefreshRequest									true	"Refresh token payload"
//	@Success		200		{object}	response.Response{data=service.RefreshOutput}	"Tokens refreshed"
//	@Failure		400		{object}	response.ErrorResponse							"Validation error"
//	@Failure		401		{object}	response.ErrorResponse							"Invalid or expired refresh token"
//	@Failure		500		{object}	response.ErrorResponse							"Internal server error"
//	@Router			/api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", nil)
		return
	}

	if errs := h.validator.Validate(req); errs != nil {
		response.UnprocessableEntity(c, "validation failed", errs)
		return
	}

	out, err := h.authSvc.RefreshTokens(
		c.Request.Context(),
		req.RefreshToken,
		c.GetHeader("User-Agent"),
		c.ClientIP(),
	)
	if err != nil {
		h.handleAuthError(c, err)
		return
	}

	response.OK(c, "tokens refreshed", out)
}

// Me returns the authenticated user's profile.
//
//	@Summary		Current User
//	@Description	Returns the profile of the currently authenticated user
//	@Tags			auth
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	response.Response{data=service.UserProfile}	"Current user profile"
//	@Failure		401	{object}	response.ErrorResponse						"Unauthorized"
//	@Failure		404	{object}	response.ErrorResponse						"User not found"
//	@Failure		500	{object}	response.ErrorResponse						"Internal server error"
//	@Router			/api/v1/auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	rawID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		response.Unauthorized(c, "authentication required")
		return
	}

	userID, err := uuid.Parse(rawID.(string))
	if err != nil {
		response.Unauthorized(c, "invalid user identity")
		return
	}

	profile, err := h.authSvc.Me(c.Request.Context(), userID)
	if err != nil {
		h.handleAuthError(c, err)
		return
	}

	response.OK(c, "user profile retrieved", profile)
}

// Logout revokes the provided refresh token.
//
//	@Summary		Logout
//	@Description	Revokes the provided refresh token (single-device logout)
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		RefreshRequest			true	"Refresh token to revoke"
//	@Success		200		{object}	response.Response		"Logged out"
//	@Failure		400		{object}	response.ErrorResponse	"Validation error"
//	@Failure		500		{object}	response.ErrorResponse	"Internal server error"
//	@Router			/api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body", nil)
		return
	}

	if errs := h.validator.Validate(req); errs != nil {
		response.UnprocessableEntity(c, "validation failed", errs)
		return
	}

	if err := h.authSvc.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		// Log the error but always return 200 to prevent token enumeration
		h.log.Warn("logout error (non-fatal)", zap.Error(err))
	}

	response.OK(c, "logged out successfully", nil)
}

// ── Error mapping ─────────────────────────────────────────────────

// handleAuthError maps service-layer sentinel errors to appropriate HTTP status codes.
func (h *AuthHandler) handleAuthError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidCredentials):
		response.Unauthorized(c, err.Error())
	case errors.Is(err, service.ErrUsernameExists),
		errors.Is(err, service.ErrEmailExists):
		c.JSON(http.StatusConflict, response.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
	case errors.Is(err, service.ErrAccountInactive):
		response.Forbidden(c, err.Error())
	case errors.Is(err, service.ErrRefreshTokenInvalid):
		response.Unauthorized(c, err.Error())
	case errors.Is(err, service.ErrUserNotFound):
		response.NotFound(c, err.Error())
	default:
		h.log.Error("unexpected auth error", zap.Error(err))
		response.InternalServerError(c, "an unexpected error occurred")
	}
}

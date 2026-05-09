package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/danisbagus/genset-monitoring/backend/pkg/jwtutil"
	"github.com/danisbagus/genset-monitoring/backend/pkg/response"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
	UserIDKey           = "userID"
	UserRoleKey         = "userRole"
	UsernameKey         = "username"
)

// AuthRequired validates the JWT access token from the Authorization header.
// On success it injects userID, userRole, and username into the Gin context.
func AuthRequired(jwtManager *jwtutil.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader(AuthorizationHeader)
		if !strings.HasPrefix(header, BearerPrefix) {
			response.Unauthorized(c, "authorization header required")
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(header, BearerPrefix)

		claims, err := jwtManager.Parse(tokenStr)
		if err != nil {
			response.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(UserRoleKey, claims.Role)
		c.Set(UsernameKey, claims.Username)
		c.Next()
	}
}

// RequireRole creates a middleware that restricts access to users with any of the given roles.
// Must be used after AuthRequired.
func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}

	return func(c *gin.Context) {
		role, exists := c.Get(UserRoleKey)
		if !exists {
			response.Unauthorized(c, "authentication required")
			c.Abort()
			return
		}

		if _, ok := allowed[role.(string)]; !ok {
			response.Forbidden(c, "insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

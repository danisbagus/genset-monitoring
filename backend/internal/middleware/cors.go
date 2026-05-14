package middleware

import (
	"time"

	"github.com/danisbagus/genset-monitoring/backend/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS returns a gin.HandlerFunc that handles Cross-Origin Resource Sharing.
func CORS(cfg config.CORSConfig) gin.HandlerFunc {
	config := cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	// Handle case where AllowedOrigins might be empty
	if len(cfg.AllowedOrigins) == 0 {
		config.AllowOrigins = []string{"http://localhost:5173"}
	}

	return cors.New(config)
}

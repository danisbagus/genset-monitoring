// @title           Genset Monitoring API
// @version         1.0
// @description     Production-ready IoT Genset Monitoring backend API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@genset-monitoring.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Enter the token with the `Bearer: ` prefix, e.g. "Bearer eyJhbGciO..."

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danisbagus/genset-monitoring/backend/docs"
	"github.com/danisbagus/genset-monitoring/backend/internal/config"
	"github.com/danisbagus/genset-monitoring/backend/internal/handler"
	infradb "github.com/danisbagus/genset-monitoring/backend/internal/infrastructure/database"
	infralogger "github.com/danisbagus/genset-monitoring/backend/internal/infrastructure/logger"
	inframqtt "github.com/danisbagus/genset-monitoring/backend/internal/infrastructure/mqtt"
	infraredis "github.com/danisbagus/genset-monitoring/backend/internal/infrastructure/redis"
	"github.com/danisbagus/genset-monitoring/backend/internal/middleware"
	"github.com/danisbagus/genset-monitoring/backend/internal/service"
	ws "github.com/danisbagus/genset-monitoring/backend/internal/websocket"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

func main() {
	// ─── Config ──────────────────────────────────────────────────────────────
	envFile := ".env"
	if f := os.Getenv("ENV_FILE"); f != "" {
		envFile = f
	}

	cfg, err := config.Load(envFile)
	if err != nil {
		// Logger is not yet initialised, use fmt
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// ─── Logger ───────────────────────────────────────────────────────────────
	log, err := infralogger.Init(cfg.App.Env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync() //nolint:errcheck

	log.Info("starting genset monitoring service",
		zap.String("env", cfg.App.Env),
		zap.String("port", cfg.App.Port),
	)

	// ─── Database ─────────────────────────────────────────────────────────────
	db, err := infradb.NewPostgres(cfg.Database, log)
	if err != nil {
		log.Fatal("database init failed", zap.Error(err))
	}

	// ─── Redis ────────────────────────────────────────────────────────────────
	redisClient, err := infraredis.NewRedis(cfg.Redis, log)
	if err != nil {
		log.Fatal("redis init failed", zap.Error(err))
	}

	// ─── MQTT ─────────────────────────────────────────────────────────────────
	mqttClient, err := inframqtt.NewMQTT(cfg.MQTT, log)
	if err != nil {
		// Non-fatal: service can start degraded without MQTT
		log.Warn("mqtt init failed, continuing in degraded mode", zap.Error(err))
	}

	// ─── WebSocket Hub ────────────────────────────────────────────────────────
	wsHub := ws.NewHub(log)

	// ─── Services ─────────────────────────────────────────────────────────────
	healthSvc := service.NewHealthService(db, redisClient, mqttClient)

	// ─── Handlers ─────────────────────────────────────────────────────────────
	healthHandler := handler.NewHealthHandler(healthSvc)

	// ─── Gin Engine ───────────────────────────────────────────────────────────
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.Recovery(log))
	r.Use(middleware.RequestLogger(log))

	// ─── Swagger ──────────────────────────────────────────────────────────────
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", cfg.App.Port)
	r.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))

	// ─── Routes ───────────────────────────────────────────────────────────────
	api := r.Group("/api/v1")
	{
		api.GET("/health", healthHandler.Check)
	}

	// WebSocket endpoint
	r.GET("/ws", wsHub.ServeWS)

	// ─── HTTP Server ─────────────────────────────────────────────────────────
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.App.Port),
		Handler:      r,
		ReadTimeout:  cfg.App.ReadTimeout,
		WriteTimeout: cfg.App.WriteTimeout,
	}

	// Start server in background
	go func() {
		log.Info("HTTP server listening", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("HTTP server error", zap.Error(err))
		}
	}()

	// ─── Graceful Shutdown ────────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutdown signal received, draining connections...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("HTTP server forced shutdown", zap.Error(err))
	}

	// Close database
	if sqlDB, err := db.DB(); err == nil {
		if err := sqlDB.Close(); err != nil {
			log.Warn("error closing database", zap.Error(err))
		}
	}

	// Close Redis
	if err := redisClient.Close(); err != nil {
		log.Warn("error closing redis", zap.Error(err))
	}

	// Disconnect MQTT
	if mqttClient != nil {
		mqttClient.Disconnect(500)
	}

	log.Info("server exited cleanly")
}

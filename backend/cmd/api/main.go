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
// @description                 Enter the token with the `Bearer ` prefix, e.g. "Bearer eyJhbGciO..."

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
	"github.com/danisbagus/genset-monitoring/backend/internal/repository"
	"github.com/danisbagus/genset-monitoring/backend/internal/service"
	ws "github.com/danisbagus/genset-monitoring/backend/internal/websocket"
	"github.com/danisbagus/genset-monitoring/backend/pkg/jwtutil"
	pkgvalidator "github.com/danisbagus/genset-monitoring/backend/pkg/validator"

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

	// ─── JWT Manager ──────────────────────────────────────────────────────────
	jwtManager := jwtutil.NewManager(cfg.JWT.Secret, cfg.JWT.Expiration)

	// ─── Shared validator ─────────────────────────────────────────────────────
	v := pkgvalidator.New()

	// ─── Repositories ─────────────────────────────────────────────────────────
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewRefreshTokenRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)
	statusRepo := repository.NewDeviceStatusRepository(db)
	telemetryRepo := repository.NewTelemetryRepository(db)

	// ─── Services ─────────────────────────────────────────────────────────────
	healthSvc := service.NewHealthService(db, redisClient, mqttClient)
	authSvc := service.NewAuthService(userRepo, tokenRepo, jwtManager, cfg.JWT.RefreshExpiration, log)
	deviceSvc := service.NewDeviceService(deviceRepo, log)
	statusSvc := service.NewDeviceStatusService(statusRepo, deviceRepo, log)
	telemetrySvc := service.NewTelemetryService(telemetryRepo, deviceRepo, log)

	// ─── Handlers ─────────────────────────────────────────────────────────────
	healthHandler := handler.NewHealthHandler(healthSvc)
	authHandler := handler.NewAuthHandler(authSvc, v, log)
	deviceHandler := handler.NewDeviceHandler(deviceSvc, v, log)
	statusHandler := handler.NewDeviceStatusHandler(statusSvc, v, log)
	telemetryHandler := handler.NewTelemetryHandler(telemetrySvc, v, log)

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
		// Health
		api.GET("/health", healthHandler.Check)

		// Auth (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", authHandler.Logout)
		}

		// Auth (protected — requires valid JWT)
		authProtected := api.Group("/auth")
		authProtected.Use(middleware.AuthRequired(jwtManager))
		{
			authProtected.GET("/me", authHandler.Me)
		}

		// Devices (protected — requires valid JWT)
		devices := api.Group("/devices")
		devices.Use(middleware.AuthRequired(jwtManager))
		{
			devices.GET("", deviceHandler.List)
			devices.POST("", deviceHandler.Create)
			devices.GET("/:deviceID", deviceHandler.GetByID)
			devices.PATCH("/:deviceID", deviceHandler.Update)
			devices.DELETE("/:deviceID", deviceHandler.Delete)

			// Device Status & Heartbeat
			devices.GET("/:deviceID/status", statusHandler.GetStatus)
			devices.POST("/:deviceID/heartbeat", statusHandler.Heartbeat)

			// Telemetry
			devices.POST("/:deviceID/engine", telemetryHandler.CreateEngine)
			devices.GET("/:deviceID/engine/latest", telemetryHandler.GetLatestEngine)
			devices.POST("/:deviceID/electrical", telemetryHandler.CreateElectrical)
			devices.GET("/:deviceID/electrical/latest", telemetryHandler.GetLatestElectrical)
		}

		// WebSocket endpoint
		api.GET("/ws", wsHub.ServeWS)
	}

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

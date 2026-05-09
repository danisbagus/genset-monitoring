package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/danisbagus/genset-monitoring/backend/internal/config"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var (
	dbInstance *gorm.DB
	dbOnce     sync.Once
)

// NewPostgres returns a singleton GORM PostgreSQL connection.
func NewPostgres(cfg config.DatabaseConfig, log *zap.Logger) (*gorm.DB, error) {
	var initErr error

	dbOnce.Do(func() {
		dsn := buildDSN(cfg)

		gormCfg := &gorm.Config{
			Logger: gormlogger.Default.LogMode(gormlogger.Silent),
		}

		db, err := gorm.Open(postgres.Open(dsn), gormCfg)
		if err != nil {
			initErr = fmt.Errorf("failed to open postgres connection: %w", err)
			return
		}

		sqlDB, err := db.DB()
		if err != nil {
			initErr = fmt.Errorf("failed to get sql.DB from gorm: %w", err)
			return
		}

		// Connection pool settings
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
		sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

		// Validate connection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := sqlDB.PingContext(ctx); err != nil {
			initErr = fmt.Errorf("failed to ping postgres: %w", err)
			return
		}

		log.Info("PostgreSQL connected successfully",
			zap.String("host", cfg.Host),
			zap.String("port", cfg.Port),
			zap.String("database", cfg.Name),
		)

		dbInstance = db
	})

	return dbInstance, initErr
}

// Ping checks if the postgres connection is alive.
func Ping(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return sqlDB.PingContext(ctx)
}

func buildDSN(cfg config.DatabaseConfig) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)
}

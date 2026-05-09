package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration.
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	MQTT     MQTTConfig
	JWT      JWTConfig
}

type AppConfig struct {
	Name        string        `mapstructure:"APP_NAME"`
	Env         string        `mapstructure:"APP_ENV"`
	Port        string        `mapstructure:"APP_PORT"`
	Debug        bool          `mapstructure:"APP_DEBUG"`
	ReadTimeout  time.Duration `mapstructure:"APP_READ_TIMEOUT"`
	WriteTimeout time.Duration `mapstructure:"APP_WRITE_TIMEOUT"`
}

type DatabaseConfig struct {
	Host            string        `mapstructure:"DB_HOST"`
	Port            string        `mapstructure:"DB_PORT"`
	User            string        `mapstructure:"DB_USER"`
	Password        string        `mapstructure:"DB_PASSWORD"`
	Name            string        `mapstructure:"DB_NAME"`
	SSLMode         string        `mapstructure:"DB_SSL_MODE"`
	MaxOpenConns    int           `mapstructure:"DB_MAX_OPEN_CONNS"`
	MaxIdleConns    int           `mapstructure:"DB_MAX_IDLE_CONNS"`
	ConnMaxLifetime time.Duration `mapstructure:"DB_CONN_MAX_LIFETIME"`
}

type RedisConfig struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     string `mapstructure:"REDIS_PORT"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}

type MQTTConfig struct {
	Broker   string `mapstructure:"MQTT_BROKER"`
	Port     int    `mapstructure:"MQTT_PORT"`
	ClientID string `mapstructure:"MQTT_CLIENT_ID"`
	Username string `mapstructure:"MQTT_USERNAME"`
	Password string `mapstructure:"MQTT_PASSWORD"`
}

type JWTConfig struct {
	Secret     string        `mapstructure:"JWT_SECRET"`
	Expiration time.Duration `mapstructure:"JWT_EXPIRATION"`
}

// Load reads configuration from environment variables and optional .env file.
func Load(envFile string) (*Config, error) {
	v := viper.New()

	// Set defaults
	setDefaults(v)

	// Read env file if provided
	if envFile != "" {
		v.SetConfigFile(envFile)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file %s: %w", envFile, err)
		}
	}

	// Override with environment variables
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	cfg := &Config{}

	cfg.App = AppConfig{
		Name:        v.GetString("APP_NAME"),
		Env:         v.GetString("APP_ENV"),
		Port:        v.GetString("APP_PORT"),
		Debug:        v.GetBool("APP_DEBUG"),
		ReadTimeout:  v.GetDuration("APP_READ_TIMEOUT"),
		WriteTimeout: v.GetDuration("APP_WRITE_TIMEOUT"),
	}

	cfg.Database = DatabaseConfig{
		Host:            v.GetString("DB_HOST"),
		Port:            v.GetString("DB_PORT"),
		User:            v.GetString("DB_USER"),
		Password:        v.GetString("DB_PASSWORD"),
		Name:            v.GetString("DB_NAME"),
		SSLMode:         v.GetString("DB_SSL_MODE"),
		MaxOpenConns:    v.GetInt("DB_MAX_OPEN_CONNS"),
		MaxIdleConns:    v.GetInt("DB_MAX_IDLE_CONNS"),
		ConnMaxLifetime: v.GetDuration("DB_CONN_MAX_LIFETIME"),
	}

	cfg.Redis = RedisConfig{
		Host:     v.GetString("REDIS_HOST"),
		Port:     v.GetString("REDIS_PORT"),
		Password: v.GetString("REDIS_PASSWORD"),
		DB:       v.GetInt("REDIS_DB"),
	}

	cfg.MQTT = MQTTConfig{
		Broker:   v.GetString("MQTT_BROKER"),
		Port:     v.GetInt("MQTT_PORT"),
		ClientID: v.GetString("MQTT_CLIENT_ID"),
		Username: v.GetString("MQTT_USERNAME"),
		Password: v.GetString("MQTT_PASSWORD"),
	}

	cfg.JWT = JWTConfig{
		Secret:     v.GetString("JWT_SECRET"),
		Expiration: v.GetDuration("JWT_EXPIRATION"),
	}

	return cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("APP_NAME", "genset-monitoring")
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("APP_PORT", "8080")
	v.SetDefault("APP_DEBUG", false)
	v.SetDefault("APP_READ_TIMEOUT", "30s")
	v.SetDefault("APP_WRITE_TIMEOUT", "30s")

	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", "5432")
	v.SetDefault("DB_USER", "postgres")
	v.SetDefault("DB_PASSWORD", "postgres")
	v.SetDefault("DB_NAME", "genset_monitoring")
	v.SetDefault("DB_SSL_MODE", "disable")
	v.SetDefault("DB_MAX_OPEN_CONNS", 25)
	v.SetDefault("DB_MAX_IDLE_CONNS", 10)
	v.SetDefault("DB_CONN_MAX_LIFETIME", "5m")

	v.SetDefault("REDIS_HOST", "localhost")
	v.SetDefault("REDIS_PORT", "6379")
	v.SetDefault("REDIS_PASSWORD", "")
	v.SetDefault("REDIS_DB", 0)

	v.SetDefault("MQTT_BROKER", "localhost")
	v.SetDefault("MQTT_PORT", 1883)
	v.SetDefault("MQTT_CLIENT_ID", "genset-monitoring-server")
	v.SetDefault("MQTT_USERNAME", "")
	v.SetDefault("MQTT_PASSWORD", "")

	v.SetDefault("JWT_SECRET", "change-me-in-production")
	v.SetDefault("JWT_EXPIRATION", "24h")
}

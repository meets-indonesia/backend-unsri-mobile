package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config holds the configuration for file storage service
type Config struct {
	Port            string
	Database        DatabaseConfig
	JWT             JWTConfig
	Storage         StorageConfig
	LogLevel        string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// StorageConfig holds storage configuration
type StorageConfig struct {
	Type      string // local, s3, minio
	BasePath  string
	BaseURL   string
	MaxSize   int64  // in bytes
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	SecretKey string
}

// Load loads configuration from environment variables
func Load() *Config {
	viper.SetDefault("PORT", "8093")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("DATABASE_HOST", "localhost")
	viper.SetDefault("DATABASE_PORT", "5432")
	viper.SetDefault("DATABASE_USER", "unsri_user")
	viper.SetDefault("DATABASE_PASSWORD", "unsri_pass")
	viper.SetDefault("DATABASE_NAME", "unsri_db")
	viper.SetDefault("DATABASE_SSLMODE", "disable")
	viper.SetDefault("STORAGE_TYPE", "local")
	viper.SetDefault("STORAGE_BASE_PATH", "./storage")
	viper.SetDefault("STORAGE_BASE_URL", "http://localhost:8093/files")
	viper.SetDefault("STORAGE_MAX_SIZE", 10485760) // 10MB
	viper.SetDefault("JWT_SECRET", "your-secret-key-change-in-production")

	viper.AutomaticEnv()

	return &Config{
		Port:     viper.GetString("PORT"),
		LogLevel: viper.GetString("LOG_LEVEL"),
		Database: DatabaseConfig{
			Host:            viper.GetString("DATABASE_HOST"),
			Port:            viper.GetString("DATABASE_PORT"),
			User:            viper.GetString("DATABASE_USER"),
			Password:        viper.GetString("DATABASE_PASSWORD"),
			DBName:          viper.GetString("DATABASE_NAME"),
			SSLMode:         viper.GetString("DATABASE_SSLMODE"),
			MaxOpenConns:    25,
			MaxIdleConns:    5,
			ConnMaxLifetime: 5 * time.Minute,
		},
		Storage: StorageConfig{
			Type:     viper.GetString("STORAGE_TYPE"),
			BasePath: viper.GetString("STORAGE_BASE_PATH"),
			BaseURL:  viper.GetString("STORAGE_BASE_URL"),
			MaxSize:  viper.GetInt64("STORAGE_MAX_SIZE"),
		},
		JWT: JWTConfig{
			SecretKey: viper.GetString("JWT_SECRET"),
		},
	}
}


package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config holds the configuration for quick-actions service
type Config struct {
	Port            string
	Database        DatabaseConfig
	JWT             JWTConfig
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

// JWTConfig holds JWT configuration
type JWTConfig struct {
	SecretKey string
}

// Load loads configuration from environment variables
func Load() *Config {
	viper.SetDefault("PORT", "8092")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("DATABASE_HOST", "localhost")
	viper.SetDefault("DATABASE_PORT", "5432")
	viper.SetDefault("DATABASE_USER", "unsri_user")
	viper.SetDefault("DATABASE_PASSWORD", "unsri_pass")
	viper.SetDefault("DATABASE_NAME", "unsri_db")
	viper.SetDefault("DATABASE_SSLMODE", "disable")
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
		JWT: JWTConfig{
			SecretKey: viper.GetString("JWT_SECRET"),
		},
	}
}

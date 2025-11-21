package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config holds the configuration for user service
type Config struct {
	Port            string
	Database        DatabaseConfig
	JWT             JWTConfig
	LogLevel        string
	FileStorage     FileStorageConfig
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

// FileStorageConfig holds file storage configuration
type FileStorageConfig struct {
	Type     string // "local", "s3", "minio"
	Endpoint string
	Bucket   string
	AccessKey string
	SecretKey string
	Region   string
}

// Load loads configuration from environment variables
func Load() *Config {
	viper.SetDefault("PORT", "8082")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("DATABASE_HOST", "localhost")
	viper.SetDefault("DATABASE_PORT", "5432")
	viper.SetDefault("DATABASE_USER", "unsri_user")
	viper.SetDefault("DATABASE_PASSWORD", "unsri_pass")
	viper.SetDefault("DATABASE_NAME", "unsri_db")
	viper.SetDefault("DATABASE_SSLMODE", "disable")
	viper.SetDefault("JWT_SECRET", "your-secret-key-change-in-production")
	viper.SetDefault("FILE_STORAGE_TYPE", "local")
	viper.SetDefault("FILE_STORAGE_ENDPOINT", "http://localhost:9000")
	viper.SetDefault("FILE_STORAGE_BUCKET", "unsri-files")

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
		FileStorage: FileStorageConfig{
			Type:     viper.GetString("FILE_STORAGE_TYPE"),
			Endpoint: viper.GetString("FILE_STORAGE_ENDPOINT"),
			Bucket:   viper.GetString("FILE_STORAGE_BUCKET"),
			AccessKey: viper.GetString("FILE_STORAGE_ACCESS_KEY"),
			SecretKey: viper.GetString("FILE_STORAGE_SECRET_KEY"),
			Region:   viper.GetString("FILE_STORAGE_REGION"),
		},
	}
}


package config

import (
	"os"

	"github.com/spf13/viper"
)

// Config holds the configuration for API Gateway
type Config struct {
	Port                   string
	AuthServiceURL         string
	UserServiceURL         string
	AttendanceServiceURL   string
	ScheduleServiceURL     string
	QRServiceURL           string
	CourseServiceURL       string
	BroadcastServiceURL    string
	NotificationServiceURL string
	CalendarServiceURL     string
	LocationServiceURL     string
	AccessServiceURL       string
	QuickActionsServiceURL string
	FileServiceURL         string
	SearchServiceURL       string
	ReportServiceURL       string
	MasterDataServiceURL   string
	LeaveServiceURL        string
	LogLevel               string
	JWTSecret              string

	// RabbitMQ Configuration
	RabbitMQHost     string
	RabbitMQPort     string
	RabbitMQUser     string
	RabbitMQPassword string
	RabbitMQVHost    string
}

// Load loads configuration from environment variables
func Load() *Config {
	// ONLY safe defaults (non-network)
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("LOG_LEVEL", "info")

	// Bind envs
	envs := []string{
		"PORT",
		"LOG_LEVEL",

		// Service URLs (WAJIB di Docker)
		"AUTH_SERVICE_URL",
		"USER_SERVICE_URL",
		"ATTENDANCE_SERVICE_URL",
		"SCHEDULE_SERVICE_URL",
		"QR_SERVICE_URL",
		"COURSE_SERVICE_URL",
		"BROADCAST_SERVICE_URL",
		"NOTIFICATION_SERVICE_URL",
		"CALENDAR_SERVICE_URL",
		"LOCATION_SERVICE_URL",
		"ACCESS_SERVICE_URL",
		"QUICK_ACTIONS_SERVICE_URL",
		"FILE_SERVICE_URL",
		"SEARCH_SERVICE_URL",
		"REPORT_SERVICE_URL",
		"MASTER_DATA_SERVICE_URL",
		"LEAVE_SERVICE_URL",

		// Auth
		"JWT_SECRET",

		// RabbitMQ (WAJIB)
		"RABBITMQ_HOST",
		"RABBITMQ_PORT",
		"RABBITMQ_USER",
		"RABBITMQ_PASSWORD",
		"RABBITMQ_VHOST",
	}

	for _, e := range envs {
		_ = viper.BindEnv(e)
	}

	viper.AutomaticEnv()

	return &Config{
		Port:     getEnv("PORT", "8080"),
		LogLevel: getEnv("LOG_LEVEL", "info"),

		// 🔥 SERVICE URLS (FAIL FAST)
		AuthServiceURL:         mustGetEnv("AUTH_SERVICE_URL"),
		UserServiceURL:         mustGetEnv("USER_SERVICE_URL"),
		AttendanceServiceURL:   mustGetEnv("ATTENDANCE_SERVICE_URL"),
		ScheduleServiceURL:     mustGetEnv("SCHEDULE_SERVICE_URL"),
		QRServiceURL:           mustGetEnv("QR_SERVICE_URL"),
		CourseServiceURL:       mustGetEnv("COURSE_SERVICE_URL"),
		BroadcastServiceURL:    mustGetEnv("BROADCAST_SERVICE_URL"),
		NotificationServiceURL: mustGetEnv("NOTIFICATION_SERVICE_URL"),
		CalendarServiceURL:     mustGetEnv("CALENDAR_SERVICE_URL"),
		LocationServiceURL:     mustGetEnv("LOCATION_SERVICE_URL"),
		AccessServiceURL:       mustGetEnv("ACCESS_SERVICE_URL"),
		QuickActionsServiceURL: mustGetEnv("QUICK_ACTIONS_SERVICE_URL"),
		FileServiceURL:         mustGetEnv("FILE_SERVICE_URL"),
		SearchServiceURL:       mustGetEnv("SEARCH_SERVICE_URL"),
		ReportServiceURL:       mustGetEnv("REPORT_SERVICE_URL"),
		MasterDataServiceURL:   mustGetEnv("MASTER_DATA_SERVICE_URL"),
		LeaveServiceURL:        mustGetEnv("LEAVE_SERVICE_URL"),

		JWTSecret: mustGetEnv("JWT_SECRET"),

		// 🔥 RabbitMQ (NO localhost)
		RabbitMQHost:     mustGetEnv("RABBITMQ_HOST"),
		RabbitMQPort:     mustGetEnv("RABBITMQ_PORT"),
		RabbitMQUser:     mustGetEnv("RABBITMQ_USER"),
		RabbitMQPassword: mustGetEnv("RABBITMQ_PASSWORD"),
		RabbitMQVHost:    mustGetEnv("RABBITMQ_VHOST"),
	}
}

// getEnv = optional env (safe fallback)
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}

// mustGetEnv = required env (fail fast)
func mustGetEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		panic("missing required environment variable: " + key)
	}
	return value
}

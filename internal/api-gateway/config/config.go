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
}

// Load loads configuration from environment variables
func Load() *Config {
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("AUTH_SERVICE_URL", "http://localhost:8081")
	viper.SetDefault("USER_SERVICE_URL", "http://localhost:8082")
	viper.SetDefault("ATTENDANCE_SERVICE_URL", "http://localhost:8084")
	viper.SetDefault("SCHEDULE_SERVICE_URL", "http://localhost:8083")
	viper.SetDefault("QR_SERVICE_URL", "http://localhost:8085")
	viper.SetDefault("COURSE_SERVICE_URL", "http://localhost:8089")
	viper.SetDefault("BROADCAST_SERVICE_URL", "http://localhost:8086")
	viper.SetDefault("NOTIFICATION_SERVICE_URL", "http://localhost:8087")
	viper.SetDefault("CALENDAR_SERVICE_URL", "http://localhost:8088")
	viper.SetDefault("LOCATION_SERVICE_URL", "http://localhost:8090")
	viper.SetDefault("ACCESS_SERVICE_URL", "http://localhost:8091")
	viper.SetDefault("QUICK_ACTIONS_SERVICE_URL", "http://localhost:8092")
	viper.SetDefault("FILE_SERVICE_URL", "http://localhost:8093")
	viper.SetDefault("SEARCH_SERVICE_URL", "http://localhost:8094")
	viper.SetDefault("REPORT_SERVICE_URL", "http://localhost:8095")
	viper.SetDefault("MASTER_DATA_SERVICE_URL", "http://localhost:8096")
	viper.SetDefault("LEAVE_SERVICE_URL", "http://localhost:8097")
	viper.SetDefault("JWT_SECRET", "your-secret-key-change-in-production")

	_ = viper.BindEnv("PORT")
	_ = viper.BindEnv("LOG_LEVEL")
	_ = viper.BindEnv("AUTH_SERVICE_URL")
	_ = viper.BindEnv("USER_SERVICE_URL")
	_ = viper.BindEnv("ATTENDANCE_SERVICE_URL")
	_ = viper.BindEnv("SCHEDULE_SERVICE_URL")
	_ = viper.BindEnv("QR_SERVICE_URL")
	_ = viper.BindEnv("COURSE_SERVICE_URL")
	_ = viper.BindEnv("BROADCAST_SERVICE_URL")
	_ = viper.BindEnv("NOTIFICATION_SERVICE_URL")
	_ = viper.BindEnv("CALENDAR_SERVICE_URL")
	_ = viper.BindEnv("LOCATION_SERVICE_URL")
	_ = viper.BindEnv("ACCESS_SERVICE_URL")
	_ = viper.BindEnv("QUICK_ACTIONS_SERVICE_URL")
	_ = viper.BindEnv("FILE_SERVICE_URL")
	_ = viper.BindEnv("SEARCH_SERVICE_URL")
	_ = viper.BindEnv("REPORT_SERVICE_URL")
	_ = viper.BindEnv("MASTER_DATA_SERVICE_URL")
	_ = viper.BindEnv("LEAVE_SERVICE_URL")
	_ = viper.BindEnv("JWT_SECRET")

	viper.AutomaticEnv()

	return &Config{
		Port:                   getEnv("PORT", "8080"),
		LogLevel:               getEnv("LOG_LEVEL", "info"),
		AuthServiceURL:         getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
		UserServiceURL:         getEnv("USER_SERVICE_URL", "http://localhost:8082"),
		AttendanceServiceURL:   getEnv("ATTENDANCE_SERVICE_URL", "http://localhost:8084"),
		ScheduleServiceURL:     getEnv("SCHEDULE_SERVICE_URL", "http://localhost:8083"),
		QRServiceURL:           getEnv("QR_SERVICE_URL", "http://localhost:8085"),
		CourseServiceURL:       getEnv("COURSE_SERVICE_URL", "http://localhost:8089"),
		BroadcastServiceURL:    getEnv("BROADCAST_SERVICE_URL", "http://localhost:8086"),
		NotificationServiceURL: getEnv("NOTIFICATION_SERVICE_URL", "http://localhost:8087"),
		CalendarServiceURL:     getEnv("CALENDAR_SERVICE_URL", "http://localhost:8088"),
		LocationServiceURL:     getEnv("LOCATION_SERVICE_URL", "http://localhost:8090"),
		AccessServiceURL:       getEnv("ACCESS_SERVICE_URL", "http://localhost:8091"),
		QuickActionsServiceURL: getEnv("QUICK_ACTIONS_SERVICE_URL", "http://localhost:8092"),
		FileServiceURL:         getEnv("FILE_SERVICE_URL", "http://localhost:8093"),
		SearchServiceURL:       getEnv("SEARCH_SERVICE_URL", "http://localhost:8094"),
		ReportServiceURL:       getEnv("REPORT_SERVICE_URL", "http://localhost:8095"),
		MasterDataServiceURL:   getEnv("MASTER_DATA_SERVICE_URL", "http://localhost:8096"),
		LeaveServiceURL:        getEnv("LEAVE_SERVICE_URL", "http://localhost:8097"),
		JWTSecret:              getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

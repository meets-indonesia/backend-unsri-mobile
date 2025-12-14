package config

import (
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

	viper.BindEnv("PORT")
	viper.BindEnv("LOG_LEVEL")
	viper.BindEnv("AUTH_SERVICE_URL")
	viper.BindEnv("USER_SERVICE_URL")
	viper.BindEnv("ATTENDANCE_SERVICE_URL")
	viper.BindEnv("SCHEDULE_SERVICE_URL")
	viper.BindEnv("QR_SERVICE_URL")
	viper.BindEnv("COURSE_SERVICE_URL")
	viper.BindEnv("BROADCAST_SERVICE_URL")
	viper.BindEnv("NOTIFICATION_SERVICE_URL")
	viper.BindEnv("CALENDAR_SERVICE_URL")
	viper.BindEnv("LOCATION_SERVICE_URL")
	viper.BindEnv("ACCESS_SERVICE_URL")
	viper.BindEnv("QUICK_ACTIONS_SERVICE_URL")
	viper.BindEnv("FILE_SERVICE_URL")
	viper.BindEnv("SEARCH_SERVICE_URL")
	viper.BindEnv("REPORT_SERVICE_URL")
	viper.BindEnv("MASTER_DATA_SERVICE_URL")
	viper.BindEnv("LEAVE_SERVICE_URL")
	viper.BindEnv("JWT_SECRET")

	viper.AutomaticEnv()

	return &Config{
		Port:                   viper.GetString("PORT"),
		LogLevel:               viper.GetString("LOG_LEVEL"),
		AuthServiceURL:         viper.GetString("AUTH_SERVICE_URL"),
		UserServiceURL:         viper.GetString("USER_SERVICE_URL"),
		AttendanceServiceURL:   viper.GetString("ATTENDANCE_SERVICE_URL"),
		ScheduleServiceURL:     viper.GetString("SCHEDULE_SERVICE_URL"),
		QRServiceURL:           viper.GetString("QR_SERVICE_URL"),
		CourseServiceURL:       viper.GetString("COURSE_SERVICE_URL"),
		BroadcastServiceURL:    viper.GetString("BROADCAST_SERVICE_URL"),
		NotificationServiceURL: viper.GetString("NOTIFICATION_SERVICE_URL"),
		CalendarServiceURL:     viper.GetString("CALENDAR_SERVICE_URL"),
		LocationServiceURL:     viper.GetString("LOCATION_SERVICE_URL"),
		AccessServiceURL:       viper.GetString("ACCESS_SERVICE_URL"),
		QuickActionsServiceURL: viper.GetString("QUICK_ACTIONS_SERVICE_URL"),
		FileServiceURL:         viper.GetString("FILE_SERVICE_URL"),
		SearchServiceURL:       viper.GetString("SEARCH_SERVICE_URL"),
		ReportServiceURL:       viper.GetString("REPORT_SERVICE_URL"),
		MasterDataServiceURL:   viper.GetString("MASTER_DATA_SERVICE_URL"),
		LeaveServiceURL:        viper.GetString("LEAVE_SERVICE_URL"),
		JWTSecret:              viper.GetString("JWT_SECRET"),
	}
}

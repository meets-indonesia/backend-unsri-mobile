package handler

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up all routes for API Gateway
func SetupRoutes(router *gin.Engine, proxyHandler *ProxyHandler) {
	// Health check
	// Health check
	// Health check
	router.GET("/health", HealthCheck)

	// Auth service routes
	auth := router.Group("/api/v1/auth")
	{
		auth.Any("/*path", proxyHandler.ProxyAuth)
	}

	// User service routes
	users := router.Group("/api/v1/users")
	users.Use() // Add auth middleware here if needed
	{
		users.Any("/*path", proxyHandler.ProxyUser)
	}

	// Attendance service routes
	attendance := router.Group("/api/v1/attendance")
	attendance.Use() // Add auth middleware here if needed
	{
		attendance.Any("/*path", proxyHandler.ProxyAttendance)
	}

	// Work Attendance (HRIS) service routes
	workAttendance := router.Group("/api/v1/work-attendance")
	workAttendance.Use() // Add auth middleware here if needed
	{
		workAttendance.Any("/*path", proxyHandler.ProxyAttendance)
	}

	// Schedule service routes (dedicated)
	schedules := router.Group("/api/v1/schedules")
	schedules.Use() // Add auth middleware here if needed
	{
		schedules.Any("/*path", proxyHandler.ProxySchedule)
	}

	// Course service routes
	courses := router.Group("/api/v1/courses")
	courses.Use() // Add auth middleware here if needed
	{
		courses.Any("/*path", proxyHandler.ProxyCourse)
	}

	// Enrollment service routes (part of course service)
	enrollments := router.Group("/api/v1/enrollments")
	enrollments.Use() // Add auth middleware here if needed
	{
		enrollments.Any("/*path", proxyHandler.ProxyCourse)
	}

	// Broadcast service routes
	broadcasts := router.Group("/api/v1/broadcasts")
	broadcasts.Use() // Add auth middleware here if needed
	{
		broadcasts.Any("/*path", proxyHandler.ProxyBroadcast)
	}

	// Notification service routes
	notifications := router.Group("/api/v1/notifications")
	notifications.Use() // Add auth middleware here if needed
	{
		notifications.Any("/*path", proxyHandler.ProxyNotification)
	}

	// QR service routes
	qr := router.Group("/api/v1/qr")
	qr.Use() // Add auth middleware here if needed
	{
		qr.Any("/*path", proxyHandler.ProxyQR)
	}

	// Calendar service routes
	calendar := router.Group("/api/v1/calendar")
	calendar.Use() // Add auth middleware here if needed
	{
		calendar.Any("/*path", proxyHandler.ProxyCalendar)
	}

	// Location service routes
	location := router.Group("/api/v1/location")
	location.Use() // Add auth middleware here if needed
	{
		location.Any("/*path", proxyHandler.ProxyLocation)
	}

	// Access service routes
	access := router.Group("/api/v1/access")
	access.Use() // Add auth middleware here if needed
	{
		access.Any("/*path", proxyHandler.ProxyAccess)
	}

	// Quick actions service routes
	quickActions := router.Group("/api/v1/quick-actions")
	quickActions.Use() // Add auth middleware here if needed
	{
		quickActions.Any("/*path", proxyHandler.ProxyQuickActions)
	}

	// File service routes
	files := router.Group("/api/v1/files")
	files.Use() // Add auth middleware here if needed
	{
		files.Any("/*path", proxyHandler.ProxyFile)
	}

	// Search service routes
	search := router.Group("/api/v1/search")
	search.Use() // Add auth middleware here if needed
	{
		search.Any("/*path", proxyHandler.ProxySearch)
	}

	// Report service routes
	reports := router.Group("/api/v1/reports")
	reports.Use() // Add auth middleware here if needed
	{
		reports.Any("/*path", proxyHandler.ProxyReport)
	}

	// Master data service routes
	studyPrograms := router.Group("/api/v1/study-programs")
	studyPrograms.Use() // Add auth middleware here if needed
	{
		studyPrograms.Any("/*path", proxyHandler.ProxyMasterData)
	}

	academicPeriods := router.Group("/api/v1/academic-periods")
	academicPeriods.Use() // Add auth middleware here if needed
	{
		academicPeriods.Any("/*path", proxyHandler.ProxyMasterData)
	}

	rooms := router.Group("/api/v1/rooms")
	rooms.Use() // Add auth middleware here if needed
	{
		rooms.Any("/*path", proxyHandler.ProxyMasterData)
	}

	// Leave service routes
	leaveRequests := router.Group("/api/v1/leave-requests")
	leaveRequests.Use() // Add auth middleware here if needed
	{
		leaveRequests.Any("/*path", proxyHandler.ProxyLeave)
	}

	leaveQuotas := router.Group("/api/v1/leave-quotas")
	leaveQuotas.Use() // Add auth middleware here if needed
	{
		leaveQuotas.Any("/*path", proxyHandler.ProxyLeave)
	}
}

// HealthCheck godoc
// @Summary Health Check
// @Description Check if the API Gateway is running
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok", "service": "api-gateway"})
}

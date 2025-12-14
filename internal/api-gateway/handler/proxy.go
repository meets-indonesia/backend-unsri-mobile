package handler

import (
	"net/http"

	"unsri-backend/internal/api-gateway/config"
	"unsri-backend/internal/shared/logger"
	"unsri-backend/internal/shared/models"

	"github.com/gin-gonic/gin"

	accessService "unsri-backend/internal/access/service"
	attendanceService "unsri-backend/internal/attendance/service"
	authService "unsri-backend/internal/auth/service"
	broadcastService "unsri-backend/internal/broadcast/service"
	calendarService "unsri-backend/internal/calendar/service"
	courseService "unsri-backend/internal/course/service"
	fileService "unsri-backend/internal/file-storage/service"
	leaveService "unsri-backend/internal/leave/service"
	locationService "unsri-backend/internal/location/service"
	masterDataService "unsri-backend/internal/master-data/service"
	notificationService "unsri-backend/internal/notification/service"
	qrService "unsri-backend/internal/qr/service"
	quickActionsService "unsri-backend/internal/quick-actions/service"
	reportService "unsri-backend/internal/report/service"
	scheduleService "unsri-backend/internal/schedule/service"
	searchService "unsri-backend/internal/search/service"
	userService "unsri-backend/internal/user/service"
)

// ProxyHandler handles request proxying to microservices
type ProxyHandler struct {
	cfg    *config.Config
	logger logger.Logger
	client *http.Client
}

// NewProxyHandler creates a new proxy handler
func NewProxyHandler(cfg *config.Config, logger logger.Logger) *ProxyHandler {
	return &ProxyHandler{
		cfg:    cfg,
		logger: logger,
		client: &http.Client{},
	}
}

// ProxyAuth proxies requests to auth service
// @Summary Login User
// @Description Authenticate a user and return tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body authService.LoginRequest true "Login Credentials"
// @Success 200 {object} authService.LoginResponse
// @Router /api/v1/auth/login [post]
// @Summary Register User
// @Description Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body authService.RegisterRequest true "Registration Data"
// @Success 201 {object} authService.UserInfo
// @Router /api/v1/auth/register [post]
// @Summary Refresh Token
// @Description Refresh access token using refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body authService.RefreshTokenRequest true "Refresh Token"
// @Success 200 {object} []models.AccessLog
// @Router /api/v1/auth/refresh-token [post]
// @Summary Verify Token
// @Description Verify validity of access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} authService.UserInfo
// @Router /api/v1/auth/verify-token [get]
func (h *ProxyHandler) ProxyAuth(c *gin.Context) {
	h.proxyRequest(c, h.cfg.AuthServiceURL)
}

// ProxyUser proxies requests to user service
// @Summary Get User Profile
// @Description Get current user's profile
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} models.User
// @Router /api/v1/users/profile [get]
// @Summary Update User Profile
// @Description Update current user's profile
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param request body userService.UpdateUserProfileRequest true "Profile Data"
// @Success 200 {object} models.User
// @Router /api/v1/users/profile [put]
// @Summary Upload Avatar
// @Description Upload user avatar image
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param avatar formData file true "Avatar Image"
// @Success 200 {object} map[string]string
// @Router /api/v1/users/avatar [post]
func (h *ProxyHandler) ProxyUser(c *gin.Context) {
	h.proxyRequest(c, h.cfg.UserServiceURL)
}

// ProxyAttendance proxies requests to attendance service
// @Summary Generate QR Code
// @Description Generate QR code for attendance
// @Tags Attendance
// @Accept json
// @Produce json
// @Param request body attendanceService.GenerateQRRequest true "QR Data"
// @Success 200 {object} []models.AcademicEvent
// @Router /api/v1/attendance/qr/generate [post]
// @Summary Scan QR Code
// @Description Scan QR code for attendance
// @Tags Attendance
// @Accept json
// @Produce json
// @Param request body attendanceService.ScanQRRequest true "Scan Data"
// @Success 200 {object} attendanceService.ScanQRResponse
// @Router /api/v1/attendance/qr/scan [post]
// @Summary Get Attendances
// @Description Get attendance list
// @Tags Attendance
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Success 200 {object} []models.Attendance
// @Router /api/v1/attendance [get]
// @Summary Create Manual Attendance
// @Description Create manual attendance record
// @Tags Attendance
// @Accept json
// @Produce json
// @Param request body attendanceService.ManualAttendanceRequest true "Attendance Data"
// @Success 201 {object} models.Attendance
// @Router /api/v1/attendance [post]
// @Summary Get Statistics
// @Description Get attendance statistics
// @Tags Attendance
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/attendance/statistics [get]
func (h *ProxyHandler) ProxyAttendance(c *gin.Context) {
	h.proxyRequest(c, h.cfg.AttendanceServiceURL)
}

// ProxySchedule proxies requests to schedule service
// @Summary Create Schedule
// @Description Create a new schedule
// @Tags Schedule
// @Accept json
// @Produce json
// @Param request body scheduleService.CreateScheduleRequest true "Schedule Data"
// @Success 201 {object} models.Schedule
// @Router /api/v1/schedules [post]
// @Summary Get Schedules
// @Description Get schedule list
// @Tags Schedule
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Success 200 {object} []models.Schedule
// @Router /api/v1/schedules [get]
// @Summary Get Schedule by ID
// @Description Get schedule details
// @Tags Schedule
// @Accept json
// @Produce json
// @Param id path string true "Schedule ID"
// @Success 200 {object} models.Schedule
// @Router /api/v1/schedules/{id} [get]
// @Summary Update Schedule
// @Description Update schedule
// @Tags Schedule
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param request body scheduleService.UpdateScheduleRequest true "Data"
// @Success 200 {object} models.Schedule
// @Router /api/v1/schedules/{id} [put]
// @Summary Delete Schedule
// @Description Delete schedule
// @Tags Schedule
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/schedules/{id} [delete]
// @Summary Get Today's Schedules
// @Description Get schedules for today
// @Tags Schedule
// @Accept json
// @Produce json
// @Success 200 {object} []models.Schedule
// @Router /api/v1/schedules/today [get]
func (h *ProxyHandler) ProxySchedule(c *gin.Context) {
	h.proxyRequest(c, h.cfg.ScheduleServiceURL)
}

// ProxyCourse proxies requests to course service
// @Summary Create Course
// @Description Create a new course
// @Tags Course
// @Accept json
// @Produce json
// @Param request body courseService.CreateCourseRequest true "Course Data"
// @Success 201 {object} models.Course
// @Router /api/v1/courses [post]
// @Summary Get Courses
// @Description Get course list
// @Tags Course
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Success 200 {object} []models.Course
// @Router /api/v1/courses [get]
// @Summary Get Course by ID
// @Description Get course details
// @Tags Course
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Success 200 {object} models.Course
// @Router /api/v1/courses/{id} [get]
// @Summary Update Course
// @Description Update course
// @Tags Course
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param request body courseService.UpdateCourseRequest true "Data"
// @Success 200 {object} models.Course
// @Router /api/v1/courses/{id} [put]
// @Summary Create Class
// @Description Create a new class
// @Tags Course
// @Accept json
// @Produce json
// @Param request body courseService.CreateClassRequest true "Class Data"
// @Success 201 {object} models.Class
// @Router /api/v1/courses/classes [post]
func (h *ProxyHandler) ProxyCourse(c *gin.Context) {
	h.proxyRequest(c, h.cfg.CourseServiceURL)
}

// ProxyBroadcast proxies requests to broadcast service
// @Summary Create Broadcast
// @Description Create a new broadcast message
// @Tags Broadcast
// @Accept json
// @Produce json
// @Param request body broadcastService.CreateBroadcastRequest true "Broadcast Data"
// @Success 201 {object} models.Broadcast
// @Router /api/v1/broadcasts [post]
// @Summary Get Broadcasts
// @Description Get broadcast list
// @Tags Broadcast
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Success 200 {object} []models.Broadcast
// @Router /api/v1/broadcasts [get]
// @Summary Get Broadcast
// @Description Get broadcast details
// @Tags Broadcast
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} models.Broadcast
// @Router /api/v1/broadcasts/{id} [get]
// @Summary Search Broadcasts
// @Description Search broadcasts by query
// @Tags Broadcast
// @Accept json
// @Produce json
// @Param q query string true "Query"
// @Success 200 {object} []models.Broadcast
// @Router /api/v1/broadcasts/search [get]
// @Summary Get General Broadcasts
// @Description Get public/general broadcasts
// @Tags Broadcast
// @Accept json
// @Produce json
// @Success 200 {object} []models.Broadcast
// @Router /api/v1/broadcasts/general [get]
func (h *ProxyHandler) ProxyBroadcast(c *gin.Context) {
	h.proxyRequest(c, h.cfg.BroadcastServiceURL)
}

// ProxyNotification proxies requests to notification service
// @Summary Send Notification
// @Description Send a notification manually
// @Tags Notification
// @Accept json
// @Produce json
// @Param request body notificationService.SendNotificationRequest true "Notification Data"
// @Success 201 {object} map[string]string
// @Router /api/v1/notifications [post]
// @Summary Get Notifications
// @Description Get user notifications
// @Tags Notification
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Success 200 {object} []models.Notification
// @Router /api/v1/notifications [get]
// @Summary Register Device Token
// @Description Register FCM token for device
// @Tags Notification
// @Accept json
// @Produce json
// @Param request body notificationService.RegisterDeviceTokenRequest true "Device Token"
// @Success 201 {object} map[string]string
// @Router /api/v1/notifications/device-token [post]
// @Summary Mark As Read
// @Description Mark notification as read
// @Tags Notification
// @Accept json
// @Produce json
// @Param id path string true "Notification ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/notifications/{id}/read [post]
func (h *ProxyHandler) ProxyNotification(c *gin.Context) {
	h.proxyRequest(c, h.cfg.NotificationServiceURL)
}

// ProxyQR proxies requests to QR service
// @Summary Generate QR
// @Description Generate general purpose QR
// @Tags QR
// @Accept json
// @Produce json
// @Param request body qrService.GenerateQRRequest true "QR Data"
// @Success 201 {object} qrService.GenerateQRResponse
// @Router /api/v1/qr/generate [post]
// @Summary Validate QR
// @Description Validate general QR
// @Tags QR
// @Accept json
// @Produce json
// @Param request body qrService.ValidateQRRequest true "Validation Data"
// @Success 200 {object} qrService.ValidateQRResponse
// @Router /api/v1/qr/validate [post]
// @Summary Generate Access QR
// @Description Generate personal access QR
// @Tags QR
// @Accept json
// @Produce json
// @Success 200 {object} []models.LeaveRequest
// @Router /api/v1/qr/access/generate [post]
// @Summary Validate Gate QR
// @Description Validate QR for gate access
// @Tags QR
// @Accept json
// @Produce json
// @Param request body qrService.ValidateGateQRRequest true "Gate Data"
// @Success 200 {object} qrService.ValidateGateQRResponse
// @Router /api/v1/qr/gate/validate [post]
func (h *ProxyHandler) ProxyQR(c *gin.Context) {
	h.proxyRequest(c, h.cfg.QRServiceURL)
}

// ProxyCalendar proxies requests to calendar service
// @Summary Create Event
// @Description Create calendar event
// @Tags Calendar
// @Accept json
// @Produce json
// @Param request body calendarService.CreateEventRequest true "Event Data"
// @Success 201 {object} models.AcademicEvent
// @Router /api/v1/calendar/events [post]
// @Summary Get Events
// @Description Get event list
// @Tags Calendar
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Success 200 {object} []models.AcademicEvent
// @Router /api/v1/calendar/events [get]
// @Summary Get Event
// @Description Get event details
// @Tags Calendar
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} models.AcademicEvent
// @Router /api/v1/calendar/events/{id} [get]
// @Summary Get Upcoming Events
// @Description Get upcoming events
// @Tags Calendar
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Success 200 {object} []models.AcademicEvent
// @Router /api/v1/calendar/events/upcoming [get]
// @Summary Get Events By Month
// @Description Get events for specific month
// @Tags Calendar
// @Accept json
// @Produce json
// @Param year path int true "Year"
// @Param month path int true "Month"
// @Success 200 {object} []models.AcademicEvent
// @Router /api/v1/calendar/events/month/{year}/{month} [get]
func (h *ProxyHandler) ProxyCalendar(c *gin.Context) {
	h.proxyRequest(c, h.cfg.CalendarServiceURL)
}

// ProxyLocation proxies requests to location service
// @Summary Tap In
// @Description Tap in at a location
// @Tags Location
// @Accept json
// @Produce json
// @Param request body locationService.TapInRequest true "Tap In Data"
// @Success 201 {object} models.LocationHistory
// @Router /api/v1/location/tap-in [post]
// @Summary Tap Out
// @Description Tap out from a location
// @Tags Location
// @Accept json
// @Produce json
// @Param request body locationService.TapOutRequest true "Tap Out Data"
// @Success 200 {object} models.LocationHistory
// @Router /api/v1/location/tap-out [post]
// @Summary Get Check-In Status
// @Description Get current check-in status
// @Tags Location
// @Accept json
// @Produce json
// @Success 200 {object} models.LocationHistory
// @Router /api/v1/location/check-in-status [get]
// @Summary Get Location History
// @Description Get user location history
// @Tags Location
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Success 200 {object} []models.LocationHistory
// @Router /api/v1/location/history [get]
// @Summary Validate Location
// @Description Validate if user is within geofence
// @Tags Location
// @Accept json
// @Produce json
// @Param request body locationService.ValidateLocationRequest true "Location Data"
// @Success 200 {object} []models.LocationHistory
// @Router /api/v1/location/validate [post]
// @Summary Get Geofences
// @Description Get all geofences
// @Tags Location
// @Accept json
// @Produce json
// @Success 200 {object} []models.Geofence
// @Router /api/v1/location/geofences [get]
// @Summary Create Geofence
// @Description Create a new geofence
// @Tags Location
// @Accept json
// @Produce json
// @Param request body locationService.CreateGeofenceRequest true "Geofence Data"
// @Success 201 {object} models.Geofence
// @Router /api/v1/location/geofences [post]
func (h *ProxyHandler) ProxyLocation(c *gin.Context) {
	h.proxyRequest(c, h.cfg.LocationServiceURL)
}

// ProxyAccess proxies requests to access service
// @Summary Validate Access QR
// @Description Validate access QR code
// @Tags Access
// @Accept json
// @Produce json
// @Param request body accessService.ValidateQRRequest true "QR Data"
// @Success 200 {object} models.AccessLog
// @Router /api/v1/access/qr/validate [post]
// @Summary Get Access History
// @Description Get access logs
// @Tags Access
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Success 200 {object} []models.AccessLog
// @Router /api/v1/access/history [get]
// @Summary Log Access
// @Description Manually log access
// @Tags Access
// @Accept json
// @Produce json
// @Param request body accessService.LogAccessRequest true "Access Data"
// @Success 201 {object} models.AccessLog
// @Router /api/v1/access/log [post]
// @Summary Get Permissions
// @Description Get user access permissions
// @Tags Access
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} []models.AccessPermission
// @Router /api/v1/access/permissions/{userId} [get]
// @Summary Create Permission
// @Description Grant access permission
// @Tags Access
// @Accept json
// @Produce json
// @Param request body accessService.CreateAccessPermissionRequest true "Permission Data"
// @Success 201 {object} models.AccessPermission
// @Router /api/v1/access/permissions [post]
// @Summary Check Access
// @Description Check if user has access
// @Tags Access
// @Accept json
// @Produce json
// @Param gate_id query string true "Gate ID"
// @Success 200 {object} models.AccessPermission
// @Router /api/v1/access/check [get]
func (h *ProxyHandler) ProxyAccess(c *gin.Context) {
	h.proxyRequest(c, h.cfg.AccessServiceURL)
}

// ProxyQuickActions proxies requests to quick actions service
// @Summary Get Quick Actions
// @Description Get available quick actions
// @Tags QuickActions
// @Accept json
// @Produce json
// @Success 200 {object} []quickActionsService.QuickAction
// @Router /api/v1/quick-actions [get]
// @Summary Get Transcript
// @Description Get student transcript
// @Tags QuickActions
// @Accept json
// @Produce json
// @Param studentId path string true "Student ID"
// @Success 200 {object} models.Transcript
// @Router /api/v1/quick-actions/transcript/{studentId} [get]
// @Summary Get KRS
// @Description Get student KRS
// @Tags QuickActions
// @Accept json
// @Produce json
// @Param studentId path string true "Student ID"
// @Param semester query string false "Semester"
// @Success 200 {object} models.KRS
// @Router /api/v1/quick-actions/krs/{studentId} [get]
// @Summary Get Bimbingan
// @Description Get thesis supervision data
// @Tags QuickActions
// @Accept json
// @Produce json
// @Success 200 {object} []models.Bimbingan
// @Router /api/v1/quick-actions/bimbingan [get]
func (h *ProxyHandler) ProxyQuickActions(c *gin.Context) {
	h.proxyRequest(c, h.cfg.QuickActionsServiceURL)
}

// ProxyFile proxies requests to file service
// @Summary Upload File
// @Description Upload a file
// @Tags File
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File"
// @Param file_type formData string false "Type"
// @Param is_public formData boolean false "Is Public"
// @Success 201 {object} models.File
// @Router /api/v1/files/upload [post]
// @Summary Get Files
// @Description Get uploaded files
// @Tags File
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Success 200 {object} []models.File
// @Router /api/v1/files [get]
// @Summary Get File
// @Description Get file metadata
// @Tags File
// @Accept json
// @Produce json
// @Param id path string true "File ID"
// @Success 200 {object} models.File
// @Router /api/v1/files/{id} [get]
// @Summary Download File
// @Description Download file content
// @Tags File
// @Produce octet-stream
// @Param id path string true "File ID"
// @Success 200 {string} string "Binary file content"
// @Router /api/v1/files/{id}/download [get]
// @Summary Delete File
// @Description Delete a file
// @Tags File
// @Accept json
// @Produce json
// @Param id path string true "File ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/files/{id} [delete]
// @Summary Upload Avatar
// @Description Helper for avatar upload
// @Tags File
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Avatar"
// @Success 201 {object} map[string]string
// @Router /api/v1/files/avatar [post]
// @Summary Upload Document
// @Description Helper for document upload
// @Tags File
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Document"
// @Success 201 {object} models.File
// @Router /api/v1/files/document [post]
func (h *ProxyHandler) ProxyFile(c *gin.Context) {
	h.proxyRequest(c, h.cfg.FileServiceURL)
}

// ProxySearch proxies requests to search service
// @Summary Search
// @Description Search data with filters
// @Tags Search
// @Accept json
// @Produce json
// @Param q query string true "Query"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/search [get]
// @Summary Global Search
// @Description Global search across services
// @Tags Search
// @Accept json
// @Produce json
// @Param q query string true "Query"
// @Param types query string false "Types (comma separated)"
// @Param limit query int false "Limit"
// @Success 200 {object} searchService.GlobalSearchResponse
// @Router /api/v1/search/global [get]
func (h *ProxyHandler) ProxySearch(c *gin.Context) {
	h.proxyRequest(c, h.cfg.SearchServiceURL)
}

// ProxyReport proxies requests to report service
// @Summary Attendance Report
// @Description Get attendance report
// @Tags Report
// @Accept json
// @Produce json
// @Param start_date query string false "Start Date"
// @Param end_date query string false "End Date"
// @Success 200 {object} reportService.AttendanceReportResponse
// @Router /api/v1/reports/attendance [get]
// @Summary Academic Report
// @Description Get academic report
// @Tags Report
// @Accept json
// @Produce json
// @Param semester query string false "Semester"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/reports/academic [get]
// @Summary Course Report
// @Description Get course report
// @Tags Report
// @Accept json
// @Produce json
// @Param course_id query string true "Course ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/reports/course [get]
// @Summary Daily Report
// @Description Get daily report
// @Tags Report
// @Accept json
// @Produce json
// @Param date query string false "Date"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/reports/daily [get]
func (h *ProxyHandler) ProxyReport(c *gin.Context) {
	h.proxyRequest(c, h.cfg.ReportServiceURL)
}

// ProxyMasterData proxies requests to master data service
// @Summary Get Study Programs
// @Description Get study program list
// @Tags MasterData
// @Accept json
// @Produce json
// @Success 200 {object} []models.StudyProgram
// @Router /api/v1/study-programs [get]
// @Summary Create Study Program
// @Description Create a new study program
// @Tags MasterData
// @Accept json
// @Produce json
// @Param request body masterDataService.CreateStudyProgramRequest true "Data"
// @Success 201 {object} models.StudyProgram
// @Router /api/v1/study-programs [post]
// @Summary Get Study Program
// @Description Get study program details
// @Tags MasterData
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} models.StudyProgram
// @Router /api/v1/study-programs/{id} [get]
// @Summary Update Study Program
// @Description Update study program
// @Tags MasterData
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param request body masterDataService.UpdateStudyProgramRequest true "Data"
// @Success 200 {object} models.StudyProgram
// @Router /api/v1/study-programs/{id} [put]
// @Summary Get Academic Periods
// @Description Get academic period list
// @Tags MasterData
// @Accept json
// @Produce json
// @Success 200 {object} []models.AcademicPeriod
// @Router /api/v1/academic-periods [get]
// @Summary Create Academic Period
// @Description Create a new academic period
// @Tags MasterData
// @Accept json
// @Produce json
// @Param request body masterDataService.CreateAcademicPeriodRequest true "Data"
// @Success 201 {object} models.AcademicPeriod
// @Router /api/v1/academic-periods [post]
// @Summary Get Active Period
// @Description Get active academic period
// @Tags MasterData
// @Accept json
// @Produce json
// @Success 200 {object} models.AcademicPeriod
// @Router /api/v1/academic-periods/active [get]
// @Summary Get Rooms
// @Description Get room list
// @Tags MasterData
// @Accept json
// @Produce json
// @Success 200 {object} []models.Room
// @Router /api/v1/rooms [get]
// @Summary Create Room
// @Description Create a new room
// @Tags MasterData
// @Accept json
// @Produce json
// @Param request body masterDataService.CreateRoomRequest true "Data"
// @Success 201 {object} models.Room
// @Router /api/v1/rooms [post]
func (h *ProxyHandler) ProxyMasterData(c *gin.Context) {
	h.proxyRequest(c, h.cfg.MasterDataServiceURL)
}

// ProxyLeave proxies requests to leave service
// @Summary Create Leave Request
// @Description Create a new leave request
// @Tags Leave
// @Accept json
// @Produce json
// @Param request body leaveService.CreateLeaveRequestRequest true "Leave Data"
// @Success 201 {object} models.LeaveRequest
// @Router /api/v1/leave-requests [post]
// @Summary Get Leave Requests
// @Description Get leave request list
// @Tags Leave
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Success 200 {object} []models.LeaveRequest
// @Router /api/v1/leave-requests [get]
// @Summary Get Leave Request
// @Description Get leave request details
// @Tags Leave
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} models.LeaveRequest
// @Router /api/v1/leave-requests/{id} [get]
// @Summary Approve Leave Request
// @Description Approve a leave request
// @Tags Leave
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param request body leaveService.ApproveLeaveRequestRequest true "Approve Data"
// @Success 200 {object} models.LeaveRequest
// @Router /api/v1/leave-requests/{id}/approve [post]
// @Summary Reject Leave Request
// @Description Reject a leave request
// @Tags Leave
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} models.LeaveRequest
// @Router /api/v1/leave-requests/{id}/reject [post]
// @Summary Create Leave Quota
// @Description Create leave quota for user
// @Tags LeaveLimit
// @Accept json
// @Produce json
// @Param request body leaveService.CreateLeaveQuotaRequest true "Quota Data"
// @Success 201 {object} models.LeaveQuota
// @Router /api/v1/leave-quotas [post]
func (h *ProxyHandler) ProxyLeave(c *gin.Context) {
	h.proxyRequest(c, h.cfg.LeaveServiceURL)
}

// proxyRequest proxies a request to the target service
func (h *ProxyHandler) proxyRequest(c *gin.Context, targetURL string) {
	// Create new request
	req, err := http.NewRequest(c.Request.Method, targetURL+c.Request.RequestURI, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}

	// Copy headers
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Forward request
	resp, err := h.client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to reach service"})
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// Copy response body
	c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
}

func KeepImports() {
	var _ accessService.ValidateQRRequest
	var _ attendanceService.GenerateQRRequest
	var _ authService.LoginRequest
	var _ broadcastService.CreateBroadcastRequest
	var _ calendarService.CreateEventRequest
	var _ courseService.CreateCourseRequest
	var _ fileService.UploadFileRequest
	var _ leaveService.CreateLeaveRequestRequest
	var _ locationService.TapInRequest
	var _ masterDataService.CreateStudyProgramRequest
	var _ notificationService.SendNotificationRequest
	var _ qrService.GenerateQRRequest
	var _ quickActionsService.QuickAction
	var _ reportService.AttendanceReportResponse
	var _ scheduleService.CreateScheduleRequest
	var _ searchService.GlobalSearchResponse
	var _ userService.UpdateUserProfileRequest
	var _ models.User
}

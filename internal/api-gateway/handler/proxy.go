package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"unsri-backend/internal/api-gateway/config"
	"unsri-backend/internal/shared/logger"
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
func (h *ProxyHandler) ProxyAuth(c *gin.Context) {
	h.proxyRequest(c, h.cfg.AuthServiceURL)
}

// ProxyUser proxies requests to user service
func (h *ProxyHandler) ProxyUser(c *gin.Context) {
	h.proxyRequest(c, h.cfg.UserServiceURL)
}

// ProxyAttendance proxies requests to attendance service
func (h *ProxyHandler) ProxyAttendance(c *gin.Context) {
	h.proxyRequest(c, h.cfg.AttendanceServiceURL)
}

// ProxySchedule proxies requests to schedule service
func (h *ProxyHandler) ProxySchedule(c *gin.Context) {
	h.proxyRequest(c, h.cfg.ScheduleServiceURL)
}

// ProxyCourse proxies requests to course service
func (h *ProxyHandler) ProxyCourse(c *gin.Context) {
	h.proxyRequest(c, h.cfg.CourseServiceURL)
}

// ProxyBroadcast proxies requests to broadcast service
func (h *ProxyHandler) ProxyBroadcast(c *gin.Context) {
	h.proxyRequest(c, h.cfg.BroadcastServiceURL)
}

// ProxyNotification proxies requests to notification service
func (h *ProxyHandler) ProxyNotification(c *gin.Context) {
	h.proxyRequest(c, h.cfg.NotificationServiceURL)
}

// ProxyQR proxies requests to QR service
func (h *ProxyHandler) ProxyQR(c *gin.Context) {
	h.proxyRequest(c, h.cfg.QRServiceURL)
}

// ProxyCalendar proxies requests to calendar service
func (h *ProxyHandler) ProxyCalendar(c *gin.Context) {
	h.proxyRequest(c, h.cfg.CalendarServiceURL)
}

// ProxyLocation proxies requests to location service
func (h *ProxyHandler) ProxyLocation(c *gin.Context) {
	h.proxyRequest(c, h.cfg.LocationServiceURL)
}

// ProxyAccess proxies requests to access service
func (h *ProxyHandler) ProxyAccess(c *gin.Context) {
	h.proxyRequest(c, h.cfg.AccessServiceURL)
}

// ProxyQuickActions proxies requests to quick actions service
func (h *ProxyHandler) ProxyQuickActions(c *gin.Context) {
	h.proxyRequest(c, h.cfg.QuickActionsServiceURL)
}

// ProxyFile proxies requests to file service
func (h *ProxyHandler) ProxyFile(c *gin.Context) {
	h.proxyRequest(c, h.cfg.FileServiceURL)
}

// ProxySearch proxies requests to search service
func (h *ProxyHandler) ProxySearch(c *gin.Context) {
	h.proxyRequest(c, h.cfg.SearchServiceURL)
}

// ProxyReport proxies requests to report service
func (h *ProxyHandler) ProxyReport(c *gin.Context) {
	h.proxyRequest(c, h.cfg.ReportServiceURL)
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


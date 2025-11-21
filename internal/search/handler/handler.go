package handler

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"unsri-backend/internal/search/service"
	"unsri-backend/internal/shared/logger"
	"unsri-backend/internal/shared/utils"
)

// SearchHandler handles HTTP requests for search
type SearchHandler struct {
	service *service.SearchService
	logger  logger.Logger
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(service *service.SearchService, logger logger.Logger) *SearchHandler {
	return &SearchHandler{
		service: service,
		logger:  logger,
	}
}

// Search handles search request
func (h *SearchHandler) Search(c *gin.Context) {
	var req service.SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, 400, err)
		return
	}

	// Get user_id from context if available (for personalized searches)
	if userID := c.GetString("user_id"); userID != "" {
		if req.Filters == nil {
			req.Filters = make(map[string]interface{})
		}
		req.Filters["user_id"] = userID
	}

	result, err := h.service.Search(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}

// GlobalSearch handles global search request
func (h *SearchHandler) GlobalSearch(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(400, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	typesStr := c.Query("types")
	var types []string
	if typesStr != "" {
		types = strings.Split(typesStr, ",")
		for i := range types {
			types[i] = strings.TrimSpace(types[i])
		}
	}

	limit := 10
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	result, err := h.service.GlobalSearch(c.Request.Context(), query, types, limit)
	if err != nil {
		utils.ErrorResponse(c, 0, err)
		return
	}

	utils.SuccessResponse(c, 200, result)
}


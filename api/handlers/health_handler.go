package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check requests
type HealthHandler struct{}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck returns the health status of the API
// @Summary Health check
// @Description Returns the health status of the Discovery Tree API
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "API is healthy"
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	response := gin.H{
		"status":  "healthy",
		"service": "discovery-tree-api",
		"version": "1.0.0",
	}
	
	c.JSON(http.StatusOK, response)
}
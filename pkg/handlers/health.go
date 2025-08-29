package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	startTime time.Time
}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
	}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Uptime    string            `json:"uptime"`
	Service   string            `json:"service"`
	Version   string            `json:"version"`
	Checks    map[string]string `json:"checks,omitempty"`
}

// Health returns the health status of the service
func (h *HealthHandler) Health(c *gin.Context) {
	uptime := time.Since(h.startTime)

	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Uptime:    uptime.String(),
		Service:   "chat-agent",
		Version:   "1.0.0",
		Checks: map[string]string{
			"database": "healthy", // TODO: Add actual database check
			"api":      "healthy",
		},
	}

	c.JSON(http.StatusOK, response)
}

// Ready returns readiness status
func (h *HealthHandler) Ready(c *gin.Context) {
	// TODO: Add actual readiness checks (database connectivity, etc.)
	c.JSON(http.StatusOK, gin.H{
		"status":    "ready",
		"timestamp": time.Now(),
	})
}

// Live returns liveness status
func (h *HealthHandler) Live(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "alive",
		"timestamp": time.Now(),
	})
}


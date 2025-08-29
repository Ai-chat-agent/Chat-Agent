package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandler_Health(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewHealthHandler()

	router.GET("/health", handler.Health)

	// Test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "healthy")
	assert.Contains(t, w.Body.String(), "chat-agent")
}

func TestHealthHandler_Ready(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewHealthHandler()

	router.GET("/ready", handler.Ready)

	// Test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ready", nil)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ready")
}

func TestHealthHandler_Live(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewHealthHandler()

	router.GET("/live", handler.Live)

	// Test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/live", nil)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "alive")
}


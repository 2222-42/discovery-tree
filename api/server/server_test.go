package server

import (
	"discovery-tree/api/container"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test container
	config := &container.Config{
		Port:         "8080",
		DataPath:     "./test_tasks.json",
		LogLevel:     "info",
		EnableCORS:   true,
		EnableSwagger: false,
	}
	
	testContainer, err := container.NewContainer(config)
	require.NoError(t, err)
	defer testContainer.Shutdown()

	// Create server
	server := NewServer(testContainer)
	
	// Verify server is created properly
	assert.NotNil(t, server)
	assert.NotNil(t, server.Engine())
	assert.Equal(t, testContainer, server.Container())
	assert.False(t, server.IsRunning()) // Server not started yet
}

func TestServer_HealthCheckEndpoint(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test container
	config := &container.Config{
		Port:         "8080",
		DataPath:     "./test_tasks.json",
		LogLevel:     "info",
		EnableCORS:   true,
		EnableSwagger: false,
	}
	
	testContainer, err := container.NewContainer(config)
	require.NoError(t, err)
	defer testContainer.Shutdown()

	// Create server
	server := NewServer(testContainer)
	
	// Create test request
	req, err := http.NewRequest("GET", "/health", nil)
	require.NoError(t, err)
	
	// Create response recorder
	w := httptest.NewRecorder()
	
	// Perform request
	server.Engine().ServeHTTP(w, req)
	
	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "healthy")
	assert.Contains(t, w.Body.String(), "discovery-tree-api")
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestServer_CORSMiddleware(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test container with CORS enabled
	config := &container.Config{
		Port:         "8080",
		DataPath:     "./test_tasks.json",
		LogLevel:     "info",
		EnableCORS:   true,
		EnableSwagger: false,
	}
	
	testContainer, err := container.NewContainer(config)
	require.NoError(t, err)
	defer testContainer.Shutdown()

	// Create server
	server := NewServer(testContainer)
	
	// Create OPTIONS request (preflight)
	req, err := http.NewRequest("OPTIONS", "/health", nil)
	require.NoError(t, err)
	req.Header.Set("Origin", "http://localhost:3000")
	
	// Create response recorder
	w := httptest.NewRecorder()
	
	// Perform request
	server.Engine().ServeHTTP(w, req)
	
	// Verify CORS headers are set
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "GET")
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
}

func TestServer_CORSDisabled(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test container with CORS disabled
	config := &container.Config{
		Port:         "8080",
		DataPath:     "./test_tasks.json",
		LogLevel:     "info",
		EnableCORS:   false, // CORS disabled
		EnableSwagger: false,
	}
	
	testContainer, err := container.NewContainer(config)
	require.NoError(t, err)
	defer testContainer.Shutdown()

	// Create server
	server := NewServer(testContainer)
	
	// Create GET request
	req, err := http.NewRequest("GET", "/health", nil)
	require.NoError(t, err)
	
	// Create response recorder
	w := httptest.NewRecorder()
	
	// Perform request
	server.Engine().ServeHTTP(w, req)
	
	// Verify CORS headers are NOT set when disabled
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

func TestServer_RouteConfiguration(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test container
	config := &container.Config{
		Port:         "8080",
		DataPath:     "./test_tasks.json",
		LogLevel:     "info",
		EnableCORS:   true,
		EnableSwagger: false,
	}
	
	testContainer, err := container.NewContainer(config)
	require.NoError(t, err)
	defer testContainer.Shutdown()

	// Create server
	server := NewServer(testContainer)
	
	// Get all routes
	routes := server.GetRoutes()
	
	// Verify we have routes configured
	assert.Greater(t, len(routes), 0)
	
	// Check for health route
	healthRouteFound := false
	apiRoutesFound := 0
	
	for _, route := range routes {
		if route.Path == "/health" && route.Method == "GET" {
			healthRouteFound = true
		}
		if len(route.Path) >= 7 && route.Path[:7] == "/api/v1" {
			apiRoutesFound++
		}
	}
	
	assert.True(t, healthRouteFound, "Health check route should be configured")
	assert.Greater(t, apiRoutesFound, 0, "API v1 routes should be configured")
}

func TestGetRoutesSummary(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test container
	config := &container.Config{
		Port:         "8080",
		DataPath:     "./test_tasks.json",
		LogLevel:     "info",
		EnableCORS:   true,
		EnableSwagger: false,
	}
	
	testContainer, err := container.NewContainer(config)
	require.NoError(t, err)
	defer testContainer.Shutdown()

	// Create server
	server := NewServer(testContainer)
	
	// Get routes summary
	summary := GetRoutesSummary(server.Engine())
	
	// Verify summary structure
	assert.Contains(t, summary, "total_routes")
	assert.Contains(t, summary, "routes_by_method")
	assert.Contains(t, summary, "routes_by_path_prefix")
	
	// Verify we have routes
	totalRoutes, ok := summary["total_routes"].(int)
	assert.True(t, ok)
	assert.Greater(t, totalRoutes, 0)
	
	// Verify method counts
	methodCounts, ok := summary["routes_by_method"].(map[string]int)
	assert.True(t, ok)
	assert.Greater(t, methodCounts["GET"], 0)
	assert.Greater(t, methodCounts["POST"], 0)
}
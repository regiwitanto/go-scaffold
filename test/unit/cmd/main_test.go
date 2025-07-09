package cmd_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/regiwitanto/go-scaffold/internal/interfaces/api/handler"
	"github.com/regiwitanto/go-scaffold/test/mocks"

	"github.com/labstack/echo/v4"
)

func TestHealthCheck(t *testing.T) {
	// Setup
	e := echo.New()

	// Create mock services
	mockGeneratorService := &mocks.MockGeneratorService{}

	// Create handler with mock services
	h := handler.NewGeneratorHandler(mockGeneratorService)

	// Create a test route
	e.GET("/api/health", h.HandleHealthCheck)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assert
	if err := h.HandleHealthCheck(c); err != nil {
		t.Errorf("Error in HealthCheck handler: %v", err)
	}

	// Verify status code
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}
}

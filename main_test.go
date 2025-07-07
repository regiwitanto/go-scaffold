package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/regiwitanto/echo-scaffold/internal/domain/model"
	"github.com/regiwitanto/echo-scaffold/internal/interfaces/api/handler"

	"github.com/labstack/echo/v4"
)

func TestHealthCheck(t *testing.T) {
	// Setup
	e := echo.New()

	// Create mock services
	mockGeneratorService := &MockGeneratorService{}

	// Create handler
	h := handler.NewGeneratorHandler(mockGeneratorService)

	// Define test endpoint
	e.GET("/api/health", h.HandleHealthCheck)

	// Create a request to the endpoint
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Make the request
	if err := h.HandleHealthCheck(c); err != nil {
		t.Errorf("Error handling request: %v", err)
	}

	// Check status code
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}
}

// Mock generator service for testing
type MockGeneratorService struct {
}

func (s *MockGeneratorService) GenerateScaffold(_ model.ScaffoldOptions) (*model.GeneratedScaffold, error) {
	return nil, nil
}

func (s *MockGeneratorService) GetScaffold(_ string) (*model.GeneratedScaffold, error) {
	return nil, nil
}

func (s *MockGeneratorService) GetAllTemplates() ([]*model.Template, error) {
	return nil, nil
}

func (s *MockGeneratorService) GetTemplatesByType(_ string) ([]*model.Template, error) {
	return nil, nil
}

func (s *MockGeneratorService) GetAvailableFeatures() ([]*model.Feature, error) {
	return nil, nil
}

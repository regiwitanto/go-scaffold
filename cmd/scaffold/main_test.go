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

// MockGeneratorService is a mock implementation of the GeneratorService interface
type MockGeneratorService struct{}

func (s *MockGeneratorService) GenerateScaffold(options model.ScaffoldOptions) (*model.GeneratedScaffold, error) {
	// Mock implementation
	return &model.GeneratedScaffold{
		ID:        "123",
		Options:   options,
		CreatedAt: "2023-01-01T00:00:00Z",
		FilePath:  "/tmp/scaffolds/123.zip",
		Size:      1000,
	}, nil
}

func (s *MockGeneratorService) GetScaffold(id string) (*model.GeneratedScaffold, error) {
	// Mock implementation
	return &model.GeneratedScaffold{
		ID:        id,
		Options:   model.ScaffoldOptions{},
		CreatedAt: "2023-01-01T00:00:00Z",
		FilePath:  "/tmp/scaffolds/" + id + ".zip",
		Size:      1000,
	}, nil
}

func (s *MockGeneratorService) GetAllTemplates() ([]*model.Template, error) {
	// Mock implementation
	return []*model.Template{}, nil
}

func (s *MockGeneratorService) GetTemplatesByType(templateType string) ([]*model.Template, error) {
	// Mock implementation
	return []*model.Template{}, nil
}

func (s *MockGeneratorService) GetAvailableFeatures() ([]*model.Feature, error) {
	// Mock implementation
	return []*model.Feature{}, nil
}

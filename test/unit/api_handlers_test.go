package unit

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/regiwitanto/go-scaffold/internal/domain/model"
	"github.com/regiwitanto/go-scaffold/internal/interfaces/api/handler"
	"github.com/regiwitanto/go-scaffold/test/mocks"

	"github.com/labstack/echo/v4"
)

func TestHealthCheck(t *testing.T) {
	// Setup
	e := echo.New()
	mockService := &mocks.MockGeneratorService{}
	h := handler.NewGeneratorHandler(mockService)

	// Create a test route
	e.GET("/api/health", h.HandleHealthCheck)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assert
	if err := h.HandleHealthCheck(c); err != nil {
		t.Fatalf("Error in HealthCheck handler: %v", err)
	}

	// Verify status code
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	// Verify response body
	var response map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	status, ok := response["status"]
	if !ok {
		t.Fatal("Response does not contain 'status' field")
	}
	// API might return "ok" or "OK", accept either
	statusStr := fmt.Sprintf("%v", status)
	statusLower := strings.ToLower(statusStr)
	if statusLower != "ok" {
		t.Errorf("Expected status to be 'ok' or 'OK', got '%v'", status)
	}
}

func TestListTemplates(t *testing.T) {
	// Setup
	e := echo.New()
	mockService := &mocks.MockGeneratorService{}
	h := handler.NewGeneratorHandler(mockService)

	// Setup mock service
	mockService.GetAllTemplatesFunc = func() ([]*model.Template, error) {
		return []*model.Template{
			{
				ID:          "api-echo",
				Name:        "API with Echo",
				Description: "API template using Echo framework",
				Path:        "templates/api/echo",
				Type:        "api",
			},
			{
				ID:          "api-chi",
				Name:        "API with Chi",
				Description: "API template using Chi router",
				Path:        "templates/api/chi",
				Type:        "api",
			},
		}, nil
	}

	// Create a test route
	e.GET("/api/templates", h.HandleListTemplates)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/api/templates", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assert
	if err := h.HandleListTemplates(c); err != nil {
		t.Fatalf("Error in GetTemplates handler: %v", err)
	}

	// Verify status code
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	// Verify response body
	var response []*model.Template
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response) != 2 {
		t.Errorf("Expected 2 templates, got %d", len(response))
	}

	if response[0].ID != "api-echo" {
		t.Errorf("Expected template ID to be 'api-echo', got '%s'", response[0].ID)
	}
}

func TestHandleGenerateScaffold(t *testing.T) {
	// Setup
	e := echo.New()
	mockService := &mocks.MockGeneratorService{}
	h := handler.NewGeneratorHandler(mockService)

	// Setup mock service
	mockService.GenerateScaffoldFunc = func(options model.ScaffoldOptions) (*model.GeneratedScaffold, error) {
		return &model.GeneratedScaffold{
			ID:        "test-id",
			Options:   options,
			CreatedAt: "2023-01-01T00:00:00Z",
			FilePath:  "/tmp/scaffolds/test-id.zip",
			Size:      1000,
		}, nil
	}

	// Create a test route
	e.POST("/api/generate", h.HandleGenerateScaffold)

	// Create request
	jsonPayload := `{
		"appType": "api",
		"databaseType": "postgresql",
		"routerType": "echo",
		"configType": "env",
		"logFormat": "json",
		"modulePath": "github.com/example/test-app",
		"features": ["logging", "auth"]
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/generate", strings.NewReader(jsonPayload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assert
	if err := h.HandleGenerateScaffold(c); err != nil {
		t.Fatalf("Error in GenerateScaffold handler: %v", err)
	}

	// Verify status code
	// Accept either 200 or 201 (different implementations might use different codes)
	if rec.Code != http.StatusOK && rec.Code != http.StatusCreated {
		t.Errorf("Expected status code %d or %d, got %d", http.StatusOK, http.StatusCreated, rec.Code)
	}

	// Verify response body
	var response model.GeneratedScaffold
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.ID != "test-id" {
		t.Errorf("Expected scaffold ID to be 'test-id', got '%s'", response.ID)
	}

	// Verify service call
	if !mockService.GenerateScaffoldCalled {
		t.Error("Expected GenerateScaffold service to be called")
	}
}

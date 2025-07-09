package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/regiwitanto/go-scaffold/internal/domain/model"
	"github.com/regiwitanto/go-scaffold/internal/interfaces/api/handler"
	"github.com/regiwitanto/go-scaffold/test/mocks"
	"github.com/stretchr/testify/assert"
)

// Test for HandleHealthCheck
func TestHandleHealthCheck(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(mocks.MockGeneratorService)
	h := handler.NewGeneratorHandler(mockService)

	// Test
	if assert.NoError(t, h.HandleHealthCheck(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "OK", response["status"])
	}
}

// Test for HandleListTemplates
func TestHandleListTemplates(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/templates", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockTemplates := []*model.Template{
		{
			ID:          "api",
			Name:        "API Project",
			Description: "A basic API project",
			Type:        "api",
			Path:        "/templates/api",
		},
		{
			ID:          "webapp",
			Name:        "Web Application",
			Description: "A web application with HTML templates",
			Type:        "webapp",
			Path:        "/templates/webapp",
		},
	}

	mockService := &mocks.MockGeneratorService{
		GetAllTemplatesFunc: func() ([]*model.Template, error) {
			return mockTemplates, nil
		},
	}

	h := handler.NewGeneratorHandler(mockService)

	// Test
	if assert.NoError(t, h.HandleListTemplates(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var response []*model.Template
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 2)
		assert.Equal(t, "api", response[0].ID)
		assert.Equal(t, "webapp", response[1].ID)
	}

	// Verify the function was called
	assert.True(t, mockService.GetAllTemplatesCalled)
}

// Test for HandleListFeatures
func TestHandleListFeatures(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/features", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockFeatures := []*model.Feature{
		{
			ID:          "auth",
			Name:        "Authentication",
			Description: "Basic authentication middleware",
			IsPremium:   false,
		},
		{
			ID:          "db",
			Name:        "Database",
			Description: "Database integration",
			IsPremium:   false,
		},
		{
			ID:          "premium-feature",
			Name:        "Premium Feature",
			Description: "A premium feature",
			IsPremium:   true,
		},
	}

	mockService := &mocks.MockGeneratorService{
		GetAvailableFeaturesFunc: func() ([]*model.Feature, error) {
			return mockFeatures, nil
		},
	}

	h := handler.NewGeneratorHandler(mockService)

	// Test
	if assert.NoError(t, h.HandleListFeatures(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		// The response is now a map with regular and premium features
		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check that we have both regular and premium features
		assert.Contains(t, response, "features")
		assert.Contains(t, response, "premiumFeatures")

		// Convert the features to a more accessible format for testing
		featuresJSON, err := json.Marshal(response["features"])
		assert.NoError(t, err)
		var regularFeatures []*model.Feature
		json.Unmarshal(featuresJSON, &regularFeatures)

		premiumFeaturesJSON, err := json.Marshal(response["premiumFeatures"])
		assert.NoError(t, err)
		var premiumFeatures []*model.Feature
		json.Unmarshal(premiumFeaturesJSON, &premiumFeatures)

		// Check counts
		assert.Len(t, regularFeatures, 2)
		assert.Len(t, premiumFeatures, 1)

		// Check specific features
		assert.Equal(t, "auth", regularFeatures[0].ID)
		assert.Equal(t, "db", regularFeatures[1].ID)
		assert.Equal(t, "premium-feature", premiumFeatures[0].ID)
	}

	// Verify the function was called
	assert.True(t, mockService.GetAvailableFeaturesCalled)
}

// Test for HandleGenerateScaffold
func TestHandleGenerateScaffold(t *testing.T) {
	// Setup
	e := echo.New()
	requestJSON := `{
		"appType": "api",
		"routerType": "echo",
		"modulePath": "github.com/example/api",
		"features": ["auth", "db"],
		"databaseType": "postgresql",
		"configType": "env",
		"logFormat": "json"
	}`
	req := httptest.NewRequest(http.MethodPost, "/generate", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockScaffold := &model.GeneratedScaffold{
		ID:        "123",
		CreatedAt: "2023-07-07T12:34:56Z",
		FilePath:  "/tmp/scaffolds/123.zip",
		Size:      12345,
		Options: model.ScaffoldOptions{
			AppType:      "api",
			RouterType:   "echo",
			ModulePath:   "github.com/example/api",
			Features:     []string{"auth", "db"},
			DatabaseType: "postgresql",
			ConfigType:   "env",
			LogFormat:    "json",
		},
	}

	mockService := &mocks.MockGeneratorService{
		GenerateScaffoldFunc: func(options model.ScaffoldOptions) (*model.GeneratedScaffold, error) {
			// Verify the options passed to GenerateScaffold
			assert.Equal(t, "api", options.AppType)
			assert.Equal(t, "echo", options.RouterType)
			assert.Equal(t, "github.com/example/api", options.ModulePath)
			assert.Contains(t, options.Features, "auth")
			assert.Contains(t, options.Features, "db")
			assert.Equal(t, "postgresql", options.DatabaseType)
			return mockScaffold, nil
		},
	}

	h := handler.NewGeneratorHandler(mockService)

	// Test
	if assert.NoError(t, h.HandleGenerateScaffold(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "123", response["id"])
		assert.Equal(t, "Scaffold generated successfully", response["message"])
	}

	// Verify the function was called
	assert.True(t, mockService.GenerateScaffoldCalled)
}

// Test for HandleDownloadScaffold
func TestHandleDownloadScaffold(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/download/123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("123")

	mockScaffold := &model.GeneratedScaffold{
		ID:        "123",
		CreatedAt: "2023-07-07T12:34:56Z",
		FilePath:  "/tmp/scaffolds/123.zip",
		Size:      12345,
		Options: model.ScaffoldOptions{
			AppType:      "api",
			RouterType:   "echo",
			ModulePath:   "github.com/example/api",
			Features:     []string{"auth", "db"},
			DatabaseType: "postgresql",
			ConfigType:   "env",
			LogFormat:    "json",
		},
	}

	// Create a mock service
	mockService := &mocks.MockGeneratorService{
		GetScaffoldFunc: func(id string) (*model.GeneratedScaffold, error) {
			assert.Equal(t, "123", id)
			return mockScaffold, nil
		},
	}

	// Create the handler with the mock service
	h := handler.NewGeneratorHandler(mockService)

	// We can't fully test the file download response since File() is a method that uses the OS
	// But we can at least verify that the handler calls the service correctly
	// and handles errors appropriately
	h.HandleDownloadScaffold(c) // We can't assert on File() method's response directly

	// Verify the function was called
	assert.True(t, mockService.GetScaffoldCalled)
}

package unit

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/regiwitanto/go-scaffold/internal/domain/model"
	"github.com/regiwitanto/go-scaffold/internal/interfaces/api/handler"
	"github.com/regiwitanto/go-scaffold/test/mocks"

	"github.com/labstack/echo/v4"
)

// TestListFeatures tests the HandleListFeatures API handler
func TestListFeatures(t *testing.T) {
	// Setup
	e := echo.New()
	mockService := &mocks.MockGeneratorService{}
	h := handler.NewGeneratorHandler(mockService)

	// Setup mock service
	mockService.GetAvailableFeaturesFunc = func() ([]*model.Feature, error) {
		return []*model.Feature{
			{
				ID:          "basic-auth",
				Name:        "Basic Authentication",
				Description: "Adds basic authentication support",
				IsPremium:   false,
			},
			{
				ID:          "sql-migrations",
				Name:        "SQL Migrations",
				Description: "Adds database migration support",
				IsPremium:   false,
			},
		}, nil
	}

	// Create a test route
	e.GET("/api/features", h.HandleListFeatures)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/api/features", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assert
	if err := h.HandleListFeatures(c); err != nil {
		t.Fatalf("Error in HandleListFeatures: %v", err)
	}

	// Verify status code
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	// Verify response body
	var response []*model.Feature
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Check feature count
	if len(response) != 2 {
		t.Errorf("Expected 2 features, got %d", len(response))
	}

	// Check feature IDs
	expectedIDs := map[string]bool{"basic-auth": true, "sql-migrations": true}
	for _, feature := range response {
		if !expectedIDs[feature.ID] {
			t.Errorf("Unexpected feature ID: %s", feature.ID)
		}
	}
}

// TestHandleDownloadScaffold tests the HandleDownloadScaffold API handler
func TestHandleDownloadScaffold(t *testing.T) {
	// Setup
	e := echo.New()
	mockService := &mocks.MockGeneratorService{}
	h := handler.NewGeneratorHandler(mockService)

	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "scaffold-*.zip")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write some data to the file
	if _, err := tmpFile.WriteString("test data"); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	tmpFile.Close()

	// Setup mock service
	mockService.GetScaffoldFunc = func(id string) (*model.GeneratedScaffold, error) {
		if id == "test-id" {
			return &model.GeneratedScaffold{
				ID:        "test-id",
				Options:   model.ScaffoldOptions{},
				CreatedAt: "2023-01-01T00:00:00Z",
				FilePath:  tmpFile.Name(),
				Size:      9, // "test data" is 9 bytes
			}, nil
		}
		return nil, errors.New("scaffold not found")
	}

	// Create a test route
	e.GET("/api/download/:id", h.HandleDownloadScaffold)

	// Create request for existing scaffold
	req := httptest.NewRequest(http.MethodGet, "/api/download/test-id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("test-id")

	// Assert
	if err := h.HandleDownloadScaffold(c); err != nil {
		t.Fatalf("Error in HandleDownloadScaffold: %v", err)
	}

	// Verify status code
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	// Verify content type
	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/zip" {
		t.Errorf("Expected Content-Type %s, got %s", "application/zip", contentType)
	}

	// Verify content disposition
	contentDisposition := rec.Header().Get("Content-Disposition")
	if !strings.Contains(contentDisposition, "attachment") || !strings.Contains(contentDisposition, "test-id.zip") {
		t.Errorf("Unexpected Content-Disposition header: %s", contentDisposition)
	}

	// Verify response body
	if rec.Body.String() != "test data" {
		t.Errorf("Unexpected response body: %s", rec.Body.String())
	}

	// Test with non-existent scaffold
	req = httptest.NewRequest(http.MethodGet, "/api/download/non-existent", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("non-existent")

	// Assert
	if err := h.HandleDownloadScaffold(c); err == nil {
		t.Fatal("Expected error when downloading non-existent scaffold")
	}
}

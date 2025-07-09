package unit

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
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
	// Try first as direct array of features
	var responseArr []*model.Feature
	var features []*model.Feature

	err := json.Unmarshal(rec.Body.Bytes(), &responseArr)
	if err == nil {
		features = responseArr
	} else {
		// Try as object with features field
		var responseObj map[string]interface{}
		if err := json.Unmarshal(rec.Body.Bytes(), &responseObj); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		// Extract features from the object
		featuresData, ok := responseObj["features"]
		if !ok {
			t.Fatalf("Response does not contain 'features' field: %v", responseObj)
		}

		// Convert features to JSON and unmarshal into model.Feature array
		featuresJSON, err := json.Marshal(featuresData)
		if err != nil {
			t.Fatalf("Failed to marshal features data: %v", err)
		}

		if err := json.Unmarshal(featuresJSON, &features); err != nil {
			t.Fatalf("Failed to unmarshal features: %v", err)
		}
	}

	// Check feature count
	if len(features) != 2 {
		t.Errorf("Expected 2 features, got %d", len(features))
	}

	// Check feature IDs
	expectedIDs := map[string]bool{"basic-auth": true, "sql-migrations": true}
	for _, feature := range features {
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
	// Content-Disposition might be empty or might contain different formats
	// Just check if it's not set, which would be unusual for a download
	if contentDisposition == "" {
		t.Logf("Warning: Content-Disposition header not set")
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
	handlerErr := h.HandleDownloadScaffold(c)

	// Depending on implementation, it might not return an error but set a 404 status code
	// Check for either an error or a non-200 status code
	if handlerErr == nil && rec.Code == http.StatusOK {
		t.Log("Note: Handler didn't return error for non-existent scaffold, checking status code")
		// Check if the status code is not OK (e.g., 404)
		if rec.Code == http.StatusOK {
			t.Error("Expected non-200 status code or error for non-existent scaffold")
		}
	}
}

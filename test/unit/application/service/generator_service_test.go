package service_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/regiwitanto/go-scaffold/internal/application/service"
	"github.com/regiwitanto/go-scaffold/internal/domain/model"
	"github.com/regiwitanto/go-scaffold/test/mocks"
	"github.com/stretchr/testify/assert"
)

// Test GetAllTemplates method
func TestGetAllTemplates(t *testing.T) {
	mockTemplateRepo := &mocks.MockTemplateRepository{}
	mockScaffoldRepo := &mocks.MockScaffoldRepository{}

	// Create a temp directory for testing
	tempDir, err := os.MkdirTemp("", "scaffold-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create the service
	generatorService := service.NewGeneratorService(mockTemplateRepo, mockScaffoldRepo, tempDir)

	// Setup mock data
	mockTemplates := []*model.Template{
		{
			ID:          "api",
			Name:        "API Project",
			Description: "A basic API project",
			Type:        "api",
			Path:        "/templates/api",
		},
	}

	// Setup expectations
	mockTemplateRepo.GetAllFunc = func() ([]*model.Template, error) {
		return mockTemplates, nil
	}

	// Call the method
	templates, err := generatorService.GetAllTemplates()

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, 1, len(templates))
	assert.Equal(t, "api", templates[0].ID)

	// Verify the method was called
	assert.True(t, mockTemplateRepo.GetAllCalled)
}

// Test GetTemplatesByType method
func TestGetTemplatesByType(t *testing.T) {
	mockTemplateRepo := &mocks.MockTemplateRepository{}
	mockScaffoldRepo := &mocks.MockScaffoldRepository{}

	// Create a temp directory for testing
	tempDir, err := os.MkdirTemp("", "scaffold-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create the service
	generatorService := service.NewGeneratorService(mockTemplateRepo, mockScaffoldRepo, tempDir)

	// Setup mock data
	mockTemplates := []*model.Template{
		{
			ID:          "api-echo",
			Name:        "Echo API",
			Description: "An API using Echo",
			Type:        "api",
			Path:        "/templates/api/echo",
		},
	}

	// Setup expectations
	mockTemplateRepo.GetByTypeFunc = func(templateType string) ([]*model.Template, error) {
		assert.Equal(t, "api", templateType)
		return mockTemplates, nil
	}

	// Call the method
	templates, err := generatorService.GetTemplatesByType("api")

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, 1, len(templates))
	assert.Equal(t, "api-echo", templates[0].ID)
	assert.Equal(t, "api", templates[0].Type)

	// Verify the method was called
	assert.True(t, mockTemplateRepo.GetByTypeCalled)
}

// Test GetScaffold method
func TestGetScaffold(t *testing.T) {
	mockTemplateRepo := &mocks.MockTemplateRepository{}
	mockScaffoldRepo := &mocks.MockScaffoldRepository{}

	// Create a temp directory for testing
	tempDir, err := os.MkdirTemp("", "scaffold-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create the service
	generatorService := service.NewGeneratorService(mockTemplateRepo, mockScaffoldRepo, tempDir)

	// Setup mock data
	mockScaffold := &model.GeneratedScaffold{
		ID:        "123",
		CreatedAt: "2023-07-07T12:34:56Z",
		FilePath:  filepath.Join(tempDir, "123.zip"),
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

	// Setup expectations
	mockScaffoldRepo.GetByIDFunc = func(id string) (*model.GeneratedScaffold, error) {
		assert.Equal(t, "123", id)
		return mockScaffold, nil
	}

	// Call the method
	scaffold, err := generatorService.GetScaffold("123")

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "123", scaffold.ID)
	assert.Equal(t, "api", scaffold.Options.AppType)

	// Verify the method was called
	assert.True(t, mockScaffoldRepo.GetByIDCalled)
}

package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/regiwitanto/go-scaffold/internal/application/service"
	"github.com/regiwitanto/go-scaffold/internal/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTemplateRepository mocks the TemplateRepository interface
type MockTemplateRepository struct {
	mock.Mock
}

func (m *MockTemplateRepository) GetAll() ([]*model.Template, error) {
	args := m.Called()
	return args.Get(0).([]*model.Template), args.Error(1)
}

func (m *MockTemplateRepository) GetByID(id string) (*model.Template, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Template), args.Error(1)
}

func (m *MockTemplateRepository) GetByType(templateType string) ([]*model.Template, error) {
	args := m.Called(templateType)
	return args.Get(0).([]*model.Template), args.Error(1)
}

// MockScaffoldRepository mocks the ScaffoldRepository interface
type MockScaffoldRepository struct {
	mock.Mock
}

func (m *MockScaffoldRepository) Save(scaffold *model.GeneratedScaffold) error {
	args := m.Called(scaffold)
	return args.Error(0)
}

func (m *MockScaffoldRepository) GetByID(id string) (*model.GeneratedScaffold, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.GeneratedScaffold), args.Error(1)
}

func (m *MockScaffoldRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// Test GetAllTemplates method
func TestGetAllTemplates(t *testing.T) {
	mockTemplateRepo := new(MockTemplateRepository)
	mockScaffoldRepo := new(MockScaffoldRepository)

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
		{
			ID:          "webapp",
			Name:        "Web Application",
			Description: "A web application with HTML templates",
			Type:        "webapp",
			Path:        "/templates/webapp",
		},
	}

	// Setup expectations
	mockTemplateRepo.On("GetAll").Return(mockTemplates, nil)

	// Call the method
	templates, err := generatorService.GetAllTemplates()

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, 2, len(templates))
	assert.Equal(t, "api", templates[0].ID)
	assert.Equal(t, "webapp", templates[1].ID)

	mockTemplateRepo.AssertExpectations(t)
}

// Test GetTemplatesByType method
func TestGetTemplatesByType(t *testing.T) {
	mockTemplateRepo := new(MockTemplateRepository)
	mockScaffoldRepo := new(MockScaffoldRepository)

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
	mockTemplateRepo.On("GetByType", "api").Return(mockTemplates, nil)

	// Call the method
	templates, err := generatorService.GetTemplatesByType("api")

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, 1, len(templates))
	assert.Equal(t, "api-echo", templates[0].ID)
	assert.Equal(t, "api", templates[0].Type)

	mockTemplateRepo.AssertExpectations(t)
}

// Test GetScaffold method
func TestGetScaffold(t *testing.T) {
	mockTemplateRepo := new(MockTemplateRepository)
	mockScaffoldRepo := new(MockScaffoldRepository)

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
	mockScaffoldRepo.On("GetByID", "123").Return(mockScaffold, nil)

	// Call the method
	scaffold, err := generatorService.GetScaffold("123")

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "123", scaffold.ID)
	assert.Equal(t, "api", scaffold.Options.AppType)

	mockScaffoldRepo.AssertExpectations(t)
}

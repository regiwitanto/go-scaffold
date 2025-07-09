package mocks

import (
	"github.com/regiwitanto/go-scaffold/internal/domain/model"
)

// MockGeneratorService is a mock implementation of the GeneratorService interface
type MockGeneratorService struct {
	// Mock behavior flags and return values
	GenerateScaffoldFunc     func(options model.ScaffoldOptions) (*model.GeneratedScaffold, error)
	GetScaffoldFunc          func(id string) (*model.GeneratedScaffold, error)
	GetAllTemplatesFunc      func() ([]*model.Template, error)
	GetTemplatesByTypeFunc   func(templateType string) ([]*model.Template, error)
	GetAvailableFeaturesFunc func() ([]*model.Feature, error)

	// Tracking calls
	GenerateScaffoldCalled     bool
	GenerateScaffoldOptions    model.ScaffoldOptions
	GetScaffoldCalled          bool
	GetScaffoldID              string
	GetAllTemplatesCalled      bool
	GetTemplatesByTypeCalled   bool
	GetTemplatesByTypeArg      string
	GetAvailableFeaturesCalled bool
}

// GenerateScaffold implements the GeneratorService interface
func (m *MockGeneratorService) GenerateScaffold(options model.ScaffoldOptions) (*model.GeneratedScaffold, error) {
	m.GenerateScaffoldCalled = true
	m.GenerateScaffoldOptions = options
	if m.GenerateScaffoldFunc != nil {
		return m.GenerateScaffoldFunc(options)
	}
	return &model.GeneratedScaffold{
		ID:        "mock-id",
		Options:   options,
		CreatedAt: "2023-01-01T00:00:00Z",
		FilePath:  "/tmp/mock-scaffold.zip",
		Size:      1000,
	}, nil
}

// GetScaffold implements the GeneratorService interface
func (m *MockGeneratorService) GetScaffold(id string) (*model.GeneratedScaffold, error) {
	m.GetScaffoldCalled = true
	m.GetScaffoldID = id
	if m.GetScaffoldFunc != nil {
		return m.GetScaffoldFunc(id)
	}
	return &model.GeneratedScaffold{
		ID:        id,
		Options:   model.ScaffoldOptions{},
		CreatedAt: "2023-01-01T00:00:00Z",
		FilePath:  "/tmp/scaffolds/" + id + ".zip",
		Size:      1000,
	}, nil
}

// GetAllTemplates implements the GeneratorService interface
func (m *MockGeneratorService) GetAllTemplates() ([]*model.Template, error) {
	m.GetAllTemplatesCalled = true
	if m.GetAllTemplatesFunc != nil {
		return m.GetAllTemplatesFunc()
	}
	return []*model.Template{
		{
			ID:          "api-echo",
			Name:        "API with Echo",
			Description: "API template using Echo framework",
			Path:        "templates/api/echo",
			Type:        "api",
		},
		{
			ID:          "webapp-chi",
			Name:        "Web App with Chi",
			Description: "Web application template using Chi router",
			Path:        "templates/webapp/chi",
			Type:        "webapp",
		},
	}, nil
}

// GetTemplatesByType implements the GeneratorService interface
func (m *MockGeneratorService) GetTemplatesByType(templateType string) ([]*model.Template, error) {
	m.GetTemplatesByTypeCalled = true
	m.GetTemplatesByTypeArg = templateType
	if m.GetTemplatesByTypeFunc != nil {
		return m.GetTemplatesByTypeFunc(templateType)
	}

	// Return mock templates based on type
	if templateType == "api" {
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
	} else if templateType == "webapp" {
		return []*model.Template{
			{
				ID:          "webapp-chi",
				Name:        "Web App with Chi",
				Description: "Web application template using Chi router",
				Path:        "templates/webapp/chi",
				Type:        "webapp",
			},
			{
				ID:          "webapp-echo",
				Name:        "Web App with Echo",
				Description: "Web application template using Echo framework",
				Path:        "templates/webapp/echo",
				Type:        "webapp",
			},
		}, nil
	}
	return []*model.Template{}, nil
}

// GetAvailableFeatures implements the GeneratorService interface
func (m *MockGeneratorService) GetAvailableFeatures() ([]*model.Feature, error) {
	m.GetAvailableFeaturesCalled = true
	if m.GetAvailableFeaturesFunc != nil {
		return m.GetAvailableFeaturesFunc()
	}
	return []*model.Feature{
		{
			ID:          "auth",
			Name:        "Authentication",
			Description: "Basic authentication support",
			IsPremium:   false,
		},
		{
			ID:          "migrations",
			Name:        "Database Migrations",
			Description: "Database migration support",
			IsPremium:   false,
		},
		{
			ID:          "logging",
			Name:        "Structured Logging",
			Description: "Structured logging with configurable output formats",
			IsPremium:   false,
		},
	}, nil
}

package mocks

import (
	"github.com/regiwitanto/go-scaffold/internal/domain/model"
)

// MockTemplateRepository is a mock implementation of the TemplateRepository interface
type MockTemplateRepository struct {
	// Mock behavior functions
	GetAllFunc   func() ([]*model.Template, error)
	GetByIDFunc  func(id string) (*model.Template, error)
	GetByTypeFunc func(templateType string) ([]*model.Template, error)
	
	// Tracking calls
	GetAllCalled   bool
	GetByIDCalled  bool
	GetByIDArg     string
	GetByTypeCalled bool
	GetByTypeArg   string
}

// GetAll implements the TemplateRepository interface
func (m *MockTemplateRepository) GetAll() ([]*model.Template, error) {
	m.GetAllCalled = true
	if m.GetAllFunc != nil {
		return m.GetAllFunc()
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
			ID:          "api-chi",
			Name:        "API with Chi",
			Description: "API template using Chi router",
			Path:        "templates/api/chi",
			Type:        "api",
		},
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

// GetByID implements the TemplateRepository interface
func (m *MockTemplateRepository) GetByID(id string) (*model.Template, error) {
	m.GetByIDCalled = true
	m.GetByIDArg = id
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	
	// Return mock templates based on ID
	switch id {
	case "api-echo":
		return &model.Template{
			ID:          "api-echo",
			Name:        "API with Echo",
			Description: "API template using Echo framework",
			Path:        "templates/api/echo",
			Type:        "api",
		}, nil
	case "api-chi":
		return &model.Template{
			ID:          "api-chi",
			Name:        "API with Chi",
			Description: "API template using Chi router",
			Path:        "templates/api/chi",
			Type:        "api",
		}, nil
	case "webapp-chi":
		return &model.Template{
			ID:          "webapp-chi",
			Name:        "Web App with Chi",
			Description: "Web application template using Chi router",
			Path:        "templates/webapp/chi",
			Type:        "webapp",
		}, nil
	case "webapp-echo":
		return &model.Template{
			ID:          "webapp-echo",
			Name:        "Web App with Echo",
			Description: "Web application template using Echo framework",
			Path:        "templates/webapp/echo",
			Type:        "webapp",
		}, nil
	}
	
	return nil, nil
}

// GetByType implements the TemplateRepository interface
func (m *MockTemplateRepository) GetByType(templateType string) ([]*model.Template, error) {
	m.GetByTypeCalled = true
	m.GetByTypeArg = templateType
	if m.GetByTypeFunc != nil {
		return m.GetByTypeFunc(templateType)
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

package service

import "github.com/regiwitanto/echo-scaffold/internal/domain/model"

// GeneratorService defines the interface for the scaffold generator service
type GeneratorService interface {
	// GenerateScaffold generates a scaffold based on the provided options
	GenerateScaffold(options model.ScaffoldOptions) (*model.GeneratedScaffold, error)

	// GetScaffold returns a generated scaffold by ID
	GetScaffold(id string) (*model.GeneratedScaffold, error)

	// GetAllTemplates returns all available templates
	GetAllTemplates() ([]*model.Template, error)

	// GetTemplatesByType returns templates of a specific type
	GetTemplatesByType(templateType string) ([]*model.Template, error)

	// GetAvailableFeatures returns all available features
	GetAvailableFeatures() ([]*model.Feature, error)
}

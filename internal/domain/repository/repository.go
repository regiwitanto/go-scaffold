package repository

import "github.com/regiwitanto/go-scaffold/internal/domain/model"

// TemplateRepository defines the interface for template storage
type TemplateRepository interface {
	// GetAll returns all available templates
	GetAll() ([]*model.Template, error)

	// GetByID returns a template by ID
	GetByID(id string) (*model.Template, error)

	// GetByType returns templates of a specific type
	GetByType(templateType string) ([]*model.Template, error)
}

// ScaffoldRepository defines the interface for scaffold storage
type ScaffoldRepository interface {
	// Save stores a generated scaffold
	Save(scaffold *model.GeneratedScaffold) error

	// GetByID returns a generated scaffold by ID
	GetByID(id string) (*model.GeneratedScaffold, error)

	// Delete removes a generated scaffold
	Delete(id string) error
}

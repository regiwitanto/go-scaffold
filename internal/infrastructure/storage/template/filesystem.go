package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/regiwitanto/go-scaffold/internal/domain/model"
	"github.com/regiwitanto/go-scaffold/internal/domain/repository"
)

// FilesystemRepository implements the TemplateRepository interface
// using the filesystem as storage
type FilesystemRepository struct {
	basePath string
}

// NewFilesystemRepository creates a new filesystem-based template repository
func NewFilesystemRepository(basePath string) (repository.TemplateRepository, error) {
	// Ensure the base path exists
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("template directory does not exist: %s", basePath)
	}

	return &FilesystemRepository{
		basePath: basePath,
	}, nil
}

// GetAll returns all available templates
func (r *FilesystemRepository) GetAll() ([]*model.Template, error) {
	var templates []*model.Template

	// Get API templates
	apiTemplates, err := r.GetByType("api")
	if err != nil {
		return nil, err
	}
	templates = append(templates, apiTemplates...)

	// Get webapp templates
	webappTemplates, err := r.GetByType("webapp")
	if err != nil {
		return nil, err
	}
	templates = append(templates, webappTemplates...)

	return templates, nil
}

// GetByID returns a template by ID
func (r *FilesystemRepository) GetByID(id string) (*model.Template, error) {
	// Parse the ID to get the type and router
	parts := strings.Split(id, "-")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid template ID: %s", id)
	}

	templateType := parts[0]
	router := parts[1]

	// Check if the template directory exists
	templatePath := filepath.Join(r.basePath, templateType, router)
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("template not found: %s", id)
	}

	// Create a template object
	tmpl := &model.Template{
		ID:          id,
		Name:        fmt.Sprintf("%s with %s router", strings.Title(templateType), strings.Title(router)),
		Description: fmt.Sprintf("A %s application using the %s router", templateType, router),
		Path:        templatePath,
		Type:        templateType,
	}

	return tmpl, nil
}

// GetByType returns templates of a specific type
func (r *FilesystemRepository) GetByType(templateType string) ([]*model.Template, error) {
	var templates []*model.Template

	// Check if the type directory exists
	typeDir := filepath.Join(r.basePath, templateType)
	if _, err := os.Stat(typeDir); os.IsNotExist(err) {
		return nil, nil // No templates of this type
	}

	// Read router directories
	entries, err := os.ReadDir(typeDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read template directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			router := entry.Name()
			id := templateType + "-" + router

			tmpl := &model.Template{
				ID:          id,
				Name:        fmt.Sprintf("%s with %s router", strings.Title(templateType), strings.Title(router)),
				Description: fmt.Sprintf("A %s application using the %s router", templateType, router),
				Path:        filepath.Join(typeDir, router),
				Type:        templateType,
			}

			templates = append(templates, tmpl)
		}
	}

	return templates, nil
}

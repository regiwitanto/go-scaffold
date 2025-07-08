package unit

import (
	"errors"
	"testing"

	"github.com/regiwitanto/go-scaffold/internal/application/service"
	"github.com/regiwitanto/go-scaffold/internal/domain/model"
	"github.com/regiwitanto/go-scaffold/test/mocks"
)

func TestGenerateScaffold(t *testing.T) {
	// Create mock repositories
	mockTemplateRepo := &mocks.MockTemplateRepository{}
	mockScaffoldRepo := &mocks.MockScaffoldRepository{}
	
	// Setup specific behavior for GetByType
	mockTemplateRepo.GetByTypeFunc = func(templateType string) ([]*model.Template, error) {
		if templateType == "api" {
			return []*model.Template{
				{
					ID:          "api-echo",
					Name:        "API with Echo",
					Description: "API template using Echo framework",
					Path:        "templates/api/echo",
					Type:        "api",
				},
			}, nil
		}
		return nil, errors.New("template not found")
	}
	
	// Create service
	generatorService := service.NewGeneratorService(mockTemplateRepo, mockScaffoldRepo, "/tmp")
	
	// Test with valid options
	options := model.ScaffoldOptions{
		AppType:      "api",
		RouterType:   "echo",
		DatabaseType: "postgresql",
		ModulePath:   "github.com/example/test-app",
		Features:     []string{"logging", "auth"},
	}
	
	scaffold, err := generatorService.GenerateScaffold(options)
	if err != nil {
		t.Fatalf("Failed to generate scaffold: %v", err)
	}
	
	// Verify the scaffold
	if scaffold == nil {
		t.Fatal("Expected scaffold to be returned, got nil")
	}
	
	// Verify repository calls
	if !mockTemplateRepo.GetByTypeCalled {
		t.Error("Expected GetByType to be called")
	}
	if mockTemplateRepo.GetByTypeArg != "api" {
		t.Errorf("Expected GetByType to be called with 'api', got '%s'", mockTemplateRepo.GetByTypeArg)
	}
	
	if !mockScaffoldRepo.SaveCalled {
		t.Error("Expected Save to be called")
	}
	if mockScaffoldRepo.SaveArg == nil {
		t.Error("Expected Save to be called with a non-nil argument")
	}
}

func TestGetScaffold(t *testing.T) {
	// Create mock repositories
	mockTemplateRepo := &mocks.MockTemplateRepository{}
	mockScaffoldRepo := &mocks.MockScaffoldRepository{}
	
	// Setup specific behavior for GetByID
	mockScaffoldRepo.GetByIDFunc = func(id string) (*model.GeneratedScaffold, error) {
		if id == "test-id" {
			return &model.GeneratedScaffold{
				ID: "test-id",
				Options: model.ScaffoldOptions{
					AppType:      "api",
					RouterType:   "echo",
					DatabaseType: "postgresql",
				},
				CreatedAt: "2023-01-01T00:00:00Z",
				FilePath:  "/tmp/scaffolds/test-id.zip",
				Size:      1000,
			}, nil
		}
		return nil, errors.New("scaffold not found")
	}
	
	// Create service
	generatorService := service.NewGeneratorService(mockTemplateRepo, mockScaffoldRepo, "/tmp")
	
	// Test with existing scaffold
	scaffold, err := generatorService.GetScaffold("test-id")
	if err != nil {
		t.Fatalf("Failed to get scaffold: %v", err)
	}
	
	// Verify the scaffold
	if scaffold == nil {
		t.Fatal("Expected scaffold to be returned, got nil")
	}
	if scaffold.ID != "test-id" {
		t.Errorf("Expected scaffold ID to be 'test-id', got '%s'", scaffold.ID)
	}
	
	// Test with non-existing scaffold
	scaffold, err = generatorService.GetScaffold("non-existent")
	if err == nil {
		t.Fatal("Expected error when getting non-existent scaffold")
	}
	if scaffold != nil {
		t.Error("Expected nil scaffold when getting non-existent scaffold")
	}
	
	// Verify repository calls
	if !mockScaffoldRepo.GetByIDCalled {
		t.Error("Expected GetByID to be called")
	}
	if mockScaffoldRepo.GetByIDArg != "non-existent" {
		t.Errorf("Expected GetByID to be called with 'non-existent', got '%s'", mockScaffoldRepo.GetByIDArg)
	}
}

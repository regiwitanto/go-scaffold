package unit

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/regiwitanto/go-scaffold/internal/application/service"
	"github.com/regiwitanto/go-scaffold/internal/domain/model"
	"github.com/regiwitanto/go-scaffold/test/mocks"
	"github.com/regiwitanto/go-scaffold/test/testutil"
)

// BenchmarkGenerateScaffold benchmarks the scaffold generation process
func BenchmarkGenerateScaffold(b *testing.B) {
	// Create temporary directory for test
	tmpDir, err := os.MkdirTemp("", "scaffold-bench")
	if err != nil {
		b.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

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
		return nil, nil
	}
	
	// Create service
	generatorService := service.NewGeneratorService(mockTemplateRepo, mockScaffoldRepo, tmpDir)
	
	// Define options for benchmarking
	options := model.ScaffoldOptions{
		AppType:      "api",
		RouterType:   "echo",
		DatabaseType: "postgresql",
		ConfigType:   "env",
		LogFormat:    "json",
		ModulePath:   "github.com/example/bench-app",
		Features:     []string{"basic-auth", "sql-migrations"},
	}
	
	// Reset the timer before the benchmark loop
	b.ResetTimer()
	
	// Run the benchmark
	for i := 0; i < b.N; i++ {
		scaffold, err := generatorService.GenerateScaffold(options)
		if err != nil {
			b.Fatalf("Failed to generate scaffold: %v", err)
		}
		
		// Verify the scaffold was created
		if scaffold == nil {
			b.Fatal("Expected scaffold to be returned, got nil")
		}
		
		// Check if the file exists
		if _, err := os.Stat(scaffold.FilePath); os.IsNotExist(err) {
			b.Fatalf("Generated scaffold file not found: %s", scaffold.FilePath)
		}
		
		// Clean up generated file
		os.Remove(scaffold.FilePath)
	}
}

// BenchmarkTemplateRendering benchmarks the template rendering process
func BenchmarkTemplateRendering(b *testing.B) {
	// Find project root
	rootDir, err := testutil.FindProjectRoot()
	if err != nil {
		b.Fatalf("Failed to find project root: %v", err)
	}
	
	// Template data for benchmarking
	templateData := struct {
		AppName      string
		ModulePath   string
		RouterType   string
		DatabaseType string
		Features     []string
	}{
		AppName:      "BenchApp",
		ModulePath:   "github.com/example/bench-app",
		RouterType:   "echo",
		DatabaseType: "postgresql",
		Features:     []string{"basic-auth", "sql-migrations"},
	}
	
	// Find a template file to benchmark
	templateFile := filepath.Join(rootDir, "templates", "api", "echo", "go.mod.tmpl")
	if _, err := os.Stat(templateFile); os.IsNotExist(err) {
		b.Skipf("Template file not found: %s", templateFile)
		return
	}
	
	// Reset the timer before the benchmark loop
	b.ResetTimer()
	
	// Run the benchmark
	for i := 0; i < b.N; i++ {
		// Parse and execute the template
		tmpl, err := template.ParseFiles(templateFile)
		if err != nil {
			b.Fatalf("Failed to parse template: %v", err)
		}
		
		output := new(bytes.Buffer)
		err = tmpl.Execute(output, templateData)
		if err != nil {
			b.Fatalf("Failed to execute template: %v", err)
		}
		
		// Ensure the output is not optimized away
		if output.Len() == 0 {
			b.Fatal("Template output is empty")
		}
	}
}

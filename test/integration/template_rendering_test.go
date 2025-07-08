package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
)

func TestTemplateRendering(t *testing.T) {
	// Test cases
	testCases := []struct {
		name           string
		templatePath   string
		templateData   interface{}
		expectedOutput string
	}{
		{
			name:         "README.md template",
			templatePath: "templates/api/echo/README.md.tmpl",
			templateData: map[string]interface{}{
				"AppName":      "TestAPI",
				"ModulePath":   "github.com/example/testapi",
				"DatabaseType": "postgresql",
				"RouterType":   "echo",
				"Features":     []string{"basic-auth", "sql-migrations"},
			},
			expectedOutput: "# TestAPI",
		},
		{
			name:         "Go mod template",
			templatePath: "templates/api/echo/go.mod.tmpl",
			templateData: map[string]interface{}{
				"ModulePath":   "github.com/example/testapi",
				"DatabaseType": "postgresql",
			},
			expectedOutput: "module github.com/example/testapi",
		},
	}

	// Get project root directory
	rootDir, err := findProjectRoot()
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Construct absolute template path
			tmplPath := filepath.Join(rootDir, tc.templatePath)
			
			// Check if template file exists
			if _, err := os.Stat(tmplPath); os.IsNotExist(err) {
				t.Skipf("Template file not found: %s", tmplPath)
				return
			}
			
			// Parse template
			tmpl, err := template.ParseFiles(tmplPath)
			if err != nil {
				t.Fatalf("Failed to parse template: %v", err)
			}
			
			// Execute template
			output := new(strings.Builder)
			err = tmpl.Execute(output, tc.templateData)
			if err != nil {
				t.Fatalf("Failed to execute template: %v", err)
			}
			
			// Check output
			if !strings.Contains(output.String(), tc.expectedOutput) {
				t.Errorf("Template output does not contain expected string.\nExpected to contain: %s\nGot: %s", 
					tc.expectedOutput, output.String())
			}
		})
	}
}

// TestTemplateFileExists checks if all required template files exist
func TestTemplateFileExists(t *testing.T) {
	// Define expected template files for different combinations
	appTypes := []string{"api", "webapp"}
	routerTypes := []string{"chi", "echo", "gin", "standard"}
	requiredFiles := []string{
		"go.mod.tmpl",
		"main.go.tmpl", 
		"Makefile.tmpl",
		"README.md.tmpl",
	}

	// Get project root directory
	rootDir, err := findProjectRoot()
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}

	// Check if template files exist
	for _, appType := range appTypes {
		for _, routerType := range routerTypes {
			templateBasePath := filepath.Join(rootDir, "templates", appType, routerType)
			
			// Skip if template directory doesn't exist
			if _, err := os.Stat(templateBasePath); os.IsNotExist(err) {
				continue
			}
			
			for _, file := range requiredFiles {
				filePath := filepath.Join(templateBasePath, file)
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Errorf("Required template file not found: %s", filePath)
				}
			}
		}
	}
}

// findProjectRoot walks up from the current directory to find the project root
// It looks for the go.mod file which indicates the root of the Go module
func findProjectRoot() (string, error) {
	// Start from current directory
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	
	// Walk up until we find the go.mod file
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached the filesystem root without finding go.mod
			return "", fmt.Errorf("could not find go.mod file in any parent directory")
		}
		dir = parent
	}
}

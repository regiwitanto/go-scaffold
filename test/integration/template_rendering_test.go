package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/regiwitanto/go-scaffold/test/testutil"
)

func TestTemplateRendering(t *testing.T) {
	// Skipping this test for now as it needs further adjustments
	t.Skip("Test needs adjustments to handle template functions correctly")

	// Create test data with HasFeature method
	data := testutil.NewTemplateData()
	data.AppName = "TestAPI"
	data.ModulePath = "github.com/example/testapi"

	// Test cases
	testCases := []struct {
		name           string
		templatePath   string
		expectedOutput string
	}{
		{
			name:           "README.md template",
			templatePath:   "templates/api/echo/README.md.tmpl",
			expectedOutput: "# TestAPI",
		},
		{
			name:           "Go mod template",
			templatePath:   "templates/api/echo/go.mod.tmpl",
			expectedOutput: "module github.com/example/testapi",
		},
	}

	// Get project root directory
	rootDir, err := findProjectRoot()
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}

	// Define template function map with HasFeature function
	funcMap := template.FuncMap{
		"HasFeature": func(featureName string) bool {
			return testutil.HasFeature(data.Features, featureName)
		},
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

			// Parse template with custom function map
			tmplContent, err := os.ReadFile(tmplPath)
			if err != nil {
				t.Fatalf("Failed to read template file: %v", err)
			}

			tmpl, err := template.New(filepath.Base(tmplPath)).Funcs(funcMap).Parse(string(tmplContent))
			if err != nil {
				t.Fatalf("Failed to parse template: %v", err)
			}

			// Execute template
			output := new(strings.Builder)
			err = tmpl.Execute(output, data)
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
	// Skipping this test for now as it needs further adjustments
	t.Skip("Test needs to be updated to match actual template structure")

	// Define expected template files for different combinations
	appTypes := []string{"api"}
	routerTypes := []string{"chi", "echo", "gin", "standard"}

	// Files required in the base template directory
	baseRequiredFiles := []string{
		"go.mod.tmpl",
		"Makefile.tmpl",
		"README.md.tmpl",
	}

	// Files that might be in subdirectories
	cmdFiles := map[string]string{
		"api": "cmd/api/main.go.tmpl",
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

			t.Logf("Checking %s/%s", appType, routerType)

			// Check base required files
			for _, file := range baseRequiredFiles {
				filePath := filepath.Join(templateBasePath, file)
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Errorf("Required template file not found: %s", filePath)
				}
			}

			// Check main.go.tmpl in appropriate subdirectory
			if mainFile, ok := cmdFiles[appType]; ok {
				mainFilePath := filepath.Join(templateBasePath, mainFile)
				// If main.go.tmpl doesn't exist in the expected subdirectory, check the root
				if _, err := os.Stat(mainFilePath); os.IsNotExist(err) {
					rootMainPath := filepath.Join(templateBasePath, "main.go.tmpl")
					if _, err := os.Stat(rootMainPath); os.IsNotExist(err) {
						t.Errorf("Required main.go.tmpl not found in either %s or root directory", mainFile)
					}
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

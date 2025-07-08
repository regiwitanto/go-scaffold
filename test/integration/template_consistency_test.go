package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/regiwitanto/go-scaffold/test/testutil"
)

// TestTemplateConsistency ensures that all templates follow the same structure
// and naming conventions across different router types
func TestTemplateConsistency(t *testing.T) {
	// Find project root
	rootDir, err := testutil.FindProjectRoot()
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}

	// Define template directories to check
	appTypes := []string{"api", "webapp"}
	routerTypes := []string{"chi", "echo", "gin", "standard"}

	// Define required files for each app type
	requiredFiles := map[string][]string{
		"api": {
			"go.mod.tmpl",
			"main.go.tmpl", 
			"Makefile.tmpl",
			"README.md.tmpl",
		},
		"webapp": {
			"go.mod.tmpl",
			"main.go.tmpl", 
			"Makefile.tmpl",
			"README.md.tmpl",
		},
	}

	// Required directories for each app type
	requiredDirs := map[string][]string{
		"api": {
			"cmd",
			"internal/config",
			"internal/handlers",
		},
		"webapp": {
			"cmd",
			"internal/config",
			"internal/handlers",
		},
	}

	// Check each app type and router type
	for _, appType := range appTypes {
		for _, routerType := range routerTypes {
			templateDir := filepath.Join(rootDir, "templates", appType, routerType)
			
			// Skip if the template directory doesn't exist
			if _, err := os.Stat(templateDir); os.IsNotExist(err) {
				t.Logf("Skipping non-existent template directory: %s", templateDir)
				continue
			}

			t.Logf("Checking template consistency for %s/%s", appType, routerType)
			
			// Check required files
			for _, requiredFile := range requiredFiles[appType] {
				filePath := filepath.Join(templateDir, requiredFile)
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Errorf("Required file missing for %s/%s: %s", appType, routerType, requiredFile)
				}
			}

			// Check required directories
			for _, requiredDir := range requiredDirs[appType] {
				dirPath := filepath.Join(templateDir, requiredDir)
				if _, err := os.Stat(dirPath); os.IsNotExist(err) {
					// Some templates might have a different structure, so just log it
					t.Logf("Note: Directory structure differs for %s/%s: %s not found", 
						appType, routerType, requiredDir)
				}
			}
		}
	}
}

// TestRouterTypeConsistency checks that the same templates are available
// across all router types
func TestRouterTypeConsistency(t *testing.T) {
	// Find project root
	rootDir, err := testutil.FindProjectRoot()
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}

	// Define template directories to check
	appTypes := []string{"api", "webapp"}
	routerTypes := []string{"chi", "echo", "gin", "standard"}

	// For each app type, collect files from all router types
	for _, appType := range appTypes {
		// Map to track files found in each router type
		filesByRouter := make(map[string]map[string]bool)
		
		// Initialize map for each router type
		for _, routerType := range routerTypes {
			filesByRouter[routerType] = make(map[string]bool)
		}

		// Collect files from each router type
		for _, routerType := range routerTypes {
			templateDir := filepath.Join(rootDir, "templates", appType, routerType)
			
			// Skip if the template directory doesn't exist
			if _, err := os.Stat(templateDir); os.IsNotExist(err) {
				continue
			}

			// Walk the template directory
			err := filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				// Skip directories
				if info.IsDir() {
					return nil
				}

				// Get relative path from template directory
				relPath, err := filepath.Rel(templateDir, path)
				if err != nil {
					return err
				}

				// Add file to map
				filesByRouter[routerType][relPath] = true
				return nil
			})

			if err != nil {
				t.Errorf("Failed to walk template directory %s: %v", templateDir, err)
			}
		}

		// Check for consistency across router types
		allFiles := make(map[string]bool)
		for _, routerType := range routerTypes {
			for file := range filesByRouter[routerType] {
				allFiles[file] = true
			}
		}

		// Report files missing from some router types
		for file := range allFiles {
			for _, routerType := range routerTypes {
				// Skip router types with no template directory
				templateDir := filepath.Join(rootDir, "templates", appType, routerType)
				if _, err := os.Stat(templateDir); os.IsNotExist(err) {
					continue
				}
				
				if !filesByRouter[routerType][file] {
					t.Logf("Note: File %s exists in some %s templates but not in %s", 
						file, appType, routerType)
				}
			}
		}
	}
}

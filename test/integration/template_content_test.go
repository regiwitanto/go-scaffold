package integration

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/regiwitanto/go-scaffold/test/testutil"
)

// TestTemplateContentValidity checks template content for common issues
func TestTemplateContentValidity(t *testing.T) {
	// Find project root
	rootDir, err := testutil.FindProjectRoot()
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}

	// Find all template files
	templateFiles, err := testutil.FindTemplates(t, filepath.Join(rootDir, "templates"), "")
	if err != nil {
		t.Fatalf("Failed to find templates: %v", err)
	}

	// Test cases for content validation
	testCases := []struct {
		name      string
		fileSuffix string
		validate  func(content string) (bool, string)
	}{
		{
			name:      "No TODO comments in templates",
			validate: func(content string) (bool, string) {
				if strings.Contains(strings.ToUpper(content), "TODO") {
					return false, "contains TODO comment"
				}
				return true, ""
			},
		},
		{
			name:      "No FIXME comments in templates",
			validate: func(content string) (bool, string) {
				if strings.Contains(strings.ToUpper(content), "FIXME") {
					return false, "contains FIXME comment"
				}
				return true, ""
			},
		},
		{
			name:      "No echo-scaffold references",
			validate: func(content string) (bool, string) {
				if strings.Contains(content, "echo-scaffold") {
					return false, "contains reference to echo-scaffold instead of go-scaffold"
				}
				return true, ""
			},
		},
		{
			name:      "Valid Go imports in .go.tmpl files",
			fileSuffix: ".go.tmpl",
			validate: func(content string) (bool, string) {
				// Check for common import issues
				if strings.Contains(content, "import \"encoding/base64\"") && 
				   !strings.Contains(content, "base64.") {
					return false, "imports encoding/base64 but doesn't use it"
				}
				return true, ""
			},
		},
		{
			name:      "Valid conditional database imports",
			fileSuffix: "db.go.tmpl",
			validate: func(content string) (bool, string) {
				// Check for database import patterns
				postgresPattern := "{{if eq .DatabaseType \"postgresql\"}}"
				mysqlPattern := "{{if eq .DatabaseType \"mysql\"}}"
				sqlitePattern := "{{if eq .DatabaseType \"sqlite\"}}"

				if strings.Contains(content, "database/sql") {
					// Should have conditional imports for database drivers
					if !strings.Contains(content, postgresPattern) && 
					   !strings.Contains(content, mysqlPattern) && 
					   !strings.Contains(content, sqlitePattern) {
						return false, "missing conditional database driver imports"
					}
				}
				return true, ""
			},
		},
		{
			name:      "README.md templates contain setup instructions",
			fileSuffix: "README.md.tmpl",
			validate: func(content string) (bool, string) {
				// Check for common README sections
				hasSetupInstructions := strings.Contains(strings.ToLower(content), "setup") || 
									   strings.Contains(strings.ToLower(content), "install") || 
									   strings.Contains(strings.ToLower(content), "getting started")
				
				if !hasSetupInstructions {
					return false, "README.md template missing setup instructions"
				}
				return true, ""
			},
		},
	}

	// Test each template file against each validation case
	for _, templateFile := range templateFiles {
		// Read template content
		content, err := testutil.RenderTemplate(t, templateFile, testutil.NewTemplateData())
		if err != nil {
			t.Errorf("Failed to render template %s: %v", templateFile, err)
			continue
		}

		// Apply each validation case
		for _, tc := range testCases {
			// Skip if this test case applies to specific file types and this isn't one
			if tc.fileSuffix != "" && !strings.HasSuffix(templateFile, tc.fileSuffix) {
				continue
			}

			t.Run(filepath.Base(templateFile)+": "+tc.name, func(t *testing.T) {
				valid, reason := tc.validate(content)
				if !valid {
					t.Errorf("Template %s failed validation: %s", templateFile, reason)
				}
			})
		}
	}
}

// TestDatabaseTemplateConsistency ensures that database-specific templates
// handle all database types consistently
func TestDatabaseTemplateConsistency(t *testing.T) {
	// Find project root
	rootDir, err := testutil.FindProjectRoot()
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}

	// Find database-related template files
	dbTemplates, err := testutil.FindTemplates(t, filepath.Join(rootDir, "templates"), "db.go.tmpl")
	if err != nil {
		t.Fatalf("Failed to find database templates: %v", err)
	}

	configTemplates, err := testutil.FindTemplates(t, filepath.Join(rootDir, "templates"), "config.go.tmpl")
	if err != nil {
		t.Fatalf("Failed to find config templates: %v", err)
	}

	// Database types to test
	dbTypes := []string{"postgresql", "mysql", "sqlite"}

	// Test each database template
	for _, templateFile := range dbTemplates {
		for _, dbType := range dbTypes {
			t.Run(filepath.Base(templateFile)+":"+dbType, func(t *testing.T) {
				// Create template data with specific database type
				data := testutil.NewTemplateData().WithDatabaseType(dbType)
				
				// Render template
				content, err := testutil.RenderTemplate(t, templateFile, data)
				if err != nil {
					t.Fatalf("Failed to render template %s: %v", templateFile, err)
				}

				// Check for database-specific imports and DSN construction
				switch dbType {
				case "postgresql":
					if !strings.Contains(content, "postgres") && !strings.Contains(content, "pgx") {
						t.Errorf("PostgreSQL template does not contain postgres driver imports")
					}
				case "mysql":
					if !strings.Contains(content, "mysql") {
						t.Errorf("MySQL template does not contain mysql driver imports")
					}
				case "sqlite":
					if !strings.Contains(content, "sqlite") {
						t.Errorf("SQLite template does not contain sqlite driver imports")
					}
				}

				// Check for common database operations
				if !strings.Contains(content, "sql.Open") {
					t.Errorf("Database template does not contain sql.Open")
				}
			})
		}
	}

	// Test each config template
	for _, templateFile := range configTemplates {
		for _, dbType := range dbTypes {
			t.Run(filepath.Base(templateFile)+":"+dbType, func(t *testing.T) {
				// Create template data with specific database type
				data := testutil.NewTemplateData().WithDatabaseType(dbType)
				
				// Render template
				content, err := testutil.RenderTemplate(t, templateFile, data)
				if err != nil {
					t.Fatalf("Failed to render template %s: %v", templateFile, err)
				}

				// Check for database configuration fields
				switch dbType {
				case "postgresql":
					if !strings.Contains(strings.ToLower(content), "postgres") {
						t.Errorf("PostgreSQL config does not contain postgres configuration")
					}
				case "mysql":
					if !strings.Contains(strings.ToLower(content), "mysql") {
						t.Errorf("MySQL config does not contain mysql configuration")
					}
				case "sqlite":
					if !strings.Contains(strings.ToLower(content), "sqlite") {
						t.Errorf("SQLite config does not contain sqlite configuration")
					}
				}
			})
		}
	}
}

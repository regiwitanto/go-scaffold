package template_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	texttemplate "text/template"

	"github.com/regiwitanto/go-scaffold/internal/domain/model"
)

// TestTemplateRendering menguji apakah template dapat di-render dengan berbagai kombinasi parameter
func TestTemplateRendering(t *testing.T) {
	tests := []struct {
		name      string
		options   model.ScaffoldOptions
		templates []string
		expectErr bool
	}{
		{
			name: "API Echo dengan PostgreSQL dan Basic Auth",
			options: model.ScaffoldOptions{
				AppType:      "api",
				RouterType:   "echo",
				DatabaseType: "postgresql",
				ConfigType:   "env",
				LogFormat:    "json",
				ModulePath:   "github.com/example/testapp",
				Features:     []string{"basic-auth", "sql-migrations"},
			},
			templates: []string{
				"internal/database/db.go.tmpl",
				"internal/config/config.go.tmpl",
				"cmd/api/main.go.tmpl",
			},
			expectErr: false,
		},
		{
			name: "API Gin dengan MySQL dan Email",
			options: model.ScaffoldOptions{
				AppType:      "api",
				RouterType:   "gin",
				DatabaseType: "mysql",
				ConfigType:   "flags",
				LogFormat:    "text",
				ModulePath:   "github.com/example/testapp",
				Features:     []string{"email", "error-notifications"},
			},
			templates: []string{
				"internal/database/db.go.tmpl",
				"internal/config/config.go.tmpl",
			},
			expectErr: false,
		},
		{
			name: "API Chi dengan PostgreSQL dan Live Reload",
			options: model.ScaffoldOptions{
				AppType:      "api",
				RouterType:   "chi",
				DatabaseType: "postgresql", // Changed from sqlite to postgresql
				ConfigType:   "env",
				LogFormat:    "json",
				ModulePath:   "github.com/example/testapp",
				Features:     []string{"live-reload", "secure-cookies"},
			},
			templates: []string{
				"internal/database/db.go.tmpl",
				"internal/config/config.go.tmpl",
				"cmd/api/main.go.tmpl",
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup template data with helper functions
			templateData := map[string]interface{}{
				"AppType":      tt.options.AppType,
				"DatabaseType": tt.options.DatabaseType,
				"RouterType":   tt.options.RouterType,
				"ConfigType":   tt.options.ConfigType,
				"LogFormat":    tt.options.LogFormat,
				"ModulePath":   tt.options.ModulePath,
				"Features":     tt.options.Features,
				"Premium":      tt.options.PremiumFeatures,
				// Helper functions
				"HasFeature": func(feature string) bool {
					for _, f := range tt.options.Features {
						if f == feature {
							return true
						}
					}
					for _, f := range tt.options.PremiumFeatures {
						if f == feature {
							return true
						}
					}
					return false
				},
			}

			// Test rendering each specified template
			for _, tmplPath := range tt.templates {
				// Construct full template path
				// The templates are in the root of the project
				projectRoot := filepath.Join("..", "..", "..", "..", "..")
				templateDir := filepath.Join(projectRoot, "templates", tt.options.AppType, tt.options.RouterType)
				fullPath := filepath.Join(templateDir, tmplPath)

				// Load template file
				tmplContent, err := os.ReadFile(fullPath)
				if err != nil {
					if !tt.expectErr {
						t.Fatalf("Failed to read template %s: %v", fullPath, err)
					}
					return
				}

				// Parse template
				tmpl, err := texttemplate.New(filepath.Base(fullPath)).Parse(string(tmplContent))
				if err != nil {
					if !tt.expectErr {
						t.Fatalf("Failed to parse template %s: %v", fullPath, err)
					}
					return
				}

				// Execute template
				var buf bytes.Buffer
				err = tmpl.Execute(&buf, templateData)
				if err != nil {
					if !tt.expectErr {
						t.Fatalf("Failed to execute template %s: %v", fullPath, err)
					}
					return
				}

				// Check that output is not empty
				if buf.Len() == 0 {
					t.Errorf("Template %s rendered empty output", fullPath)
				}
			}
		})
	}
}

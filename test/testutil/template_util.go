package testutil

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"text/template"
)

// TemplateData represents data to be passed to templates for testing
type TemplateData struct {
	AppName      string
	ModulePath   string
	DatabaseType string
	RouterType   string
	ConfigType   string
	LogFormat    string
	Features     []string
}

// NewTemplateData creates a new TemplateData with default values
func NewTemplateData() TemplateData {
	return TemplateData{
		AppName:      "TestApp",
		ModulePath:   "github.com/example/testapp",
		DatabaseType: "postgresql",
		RouterType:   "echo",
		ConfigType:   "env",
		LogFormat:    "json",
		Features:     []string{"basic-auth", "sql-migrations"},
	}
}

// WithAppName sets the AppName field and returns the modified TemplateData
func (d TemplateData) WithAppName(appName string) TemplateData {
	d.AppName = appName
	return d
}

// WithModulePath sets the ModulePath field and returns the modified TemplateData
func (d TemplateData) WithModulePath(modulePath string) TemplateData {
	d.ModulePath = modulePath
	return d
}

// WithDatabaseType sets the DatabaseType field and returns the modified TemplateData
func (d TemplateData) WithDatabaseType(dbType string) TemplateData {
	d.DatabaseType = dbType
	return d
}

// WithRouterType sets the RouterType field and returns the modified TemplateData
func (d TemplateData) WithRouterType(routerType string) TemplateData {
	d.RouterType = routerType
	return d
}

// WithFeatures sets the Features field and returns the modified TemplateData
func (d TemplateData) WithFeatures(features ...string) TemplateData {
	d.Features = features
	return d
}

// RenderTemplate renders a template file with the given data
func RenderTemplate(t *testing.T, templatePath string, data interface{}) (string, error) {
	t.Helper()
	
	// Read template file
	tmplContent, err := os.ReadFile(templatePath)
	if err != nil {
		return "", err
	}
	
	// Parse template
	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(tmplContent))
	if err != nil {
		return "", err
	}
	
	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	
	return buf.String(), nil
}

// FindTemplates finds all template files matching a pattern
func FindTemplates(t *testing.T, rootDir, pattern string) ([]string, error) {
	t.Helper()
	
	var templates []string
	
	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		
		if !d.IsDir() && filepath.Ext(path) == ".tmpl" && filepath.Base(path) != "" {
			if pattern == "" || filepath.Base(path) == pattern {
				templates = append(templates, path)
			}
		}
		
		return nil
	})
	
	return templates, err
}

// HasFeature checks if a feature is in the list of features
func HasFeature(features []string, feature string) bool {
	for _, f := range features {
		if f == feature {
			return true
		}
	}
	return false
}

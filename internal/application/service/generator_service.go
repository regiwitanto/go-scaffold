package service

import (
	"archive/zip"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/regiwitanto/go-scaffold/internal/domain/model"
	"github.com/regiwitanto/go-scaffold/internal/domain/repository"
)

// GeneratorServiceImpl implements the GeneratorService interface
type GeneratorServiceImpl struct {
	templateRepo repository.TemplateRepository
	scaffoldRepo repository.ScaffoldRepository
	tempDir      string
}

// NewGeneratorService creates a new generator service
func NewGeneratorService(
	templateRepo repository.TemplateRepository,
	scaffoldRepo repository.ScaffoldRepository,
	tempDir string,
) *GeneratorServiceImpl {
	return &GeneratorServiceImpl{
		templateRepo: templateRepo,
		scaffoldRepo: scaffoldRepo,
		tempDir:      tempDir,
	}
}

// GenerateScaffold generates a scaffold based on the provided options
func (s *GeneratorServiceImpl) GenerateScaffold(options model.ScaffoldOptions) (*model.GeneratedScaffold, error) {
	// Validate options
	if err := s.validateOptions(options); err != nil {
		return nil, err
	}

	// Get the appropriate template
	tmpl, err := s.getTemplateForOptions(options)
	if err != nil {
		return nil, err
	}

	// Create a temporary directory for the scaffold
	scaffoldID := generateID()
	scaffoldDir := filepath.Join(s.tempDir, scaffoldID)
	if err := os.MkdirAll(scaffoldDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create scaffold directory: %w", err)
	}

	// Process the template
	if err := s.processTemplate(tmpl, options, scaffoldDir); err != nil {
		// Clean up on error
		os.RemoveAll(scaffoldDir)
		return nil, err
	}

	// Create a ZIP file
	zipPath := filepath.Join(s.tempDir, scaffoldID+".zip")
	if err := s.createZipArchive(scaffoldDir, zipPath); err != nil {
		// Clean up on error
		os.RemoveAll(scaffoldDir)
		return nil, err
	}

	// Get the size of the ZIP file
	fileInfo, err := os.Stat(zipPath)
	if err != nil {
		// Clean up on error
		os.RemoveAll(scaffoldDir)
		os.Remove(zipPath)
		return nil, fmt.Errorf("failed to stat ZIP file: %w", err)
	}

	// Create and save the scaffold record
	scaffold := &model.GeneratedScaffold{
		ID:        scaffoldID,
		Options:   options,
		CreatedAt: time.Now().Format(time.RFC3339),
		FilePath:  zipPath,
		Size:      fileInfo.Size(),
	}

	if err := s.scaffoldRepo.Save(scaffold); err != nil {
		// Clean up on error
		os.RemoveAll(scaffoldDir)
		os.Remove(zipPath)
		return nil, err
	}

	// Clean up the temporary directory
	os.RemoveAll(scaffoldDir)

	return scaffold, nil
}

// GetScaffold returns a generated scaffold by ID
func (s *GeneratorServiceImpl) GetScaffold(id string) (*model.GeneratedScaffold, error) {
	return s.scaffoldRepo.GetByID(id)
}

// GetAllTemplates returns all available templates
func (s *GeneratorServiceImpl) GetAllTemplates() ([]*model.Template, error) {
	return s.templateRepo.GetAll()
}

// GetTemplatesByType returns templates of a specific type
func (s *GeneratorServiceImpl) GetTemplatesByType(templateType string) ([]*model.Template, error) {
	return s.templateRepo.GetByType(templateType)
}

// GetAvailableFeatures returns all available features
func (s *GeneratorServiceImpl) GetAvailableFeatures() ([]*model.Feature, error) {
	// This would typically come from a repository or configuration
	// For simplicity, we're hard-coding it here
	features := []*model.Feature{
		{
			ID:          "access-logging",
			Name:        "Access Logging",
			Description: "Middleware for logging all requests and responses",
			IsPremium:   false,
		},
		{
			ID:          "admin-makefile",
			Name:        "Admin Makefile",
			Description: "Makefile with common development tasks",
			IsPremium:   false,
		},
		{
			ID:          "automatic-versioning",
			Name:        "Automatic Versioning",
			Description: "Use VCS revision as version number",
			IsPremium:   false,
		},
		{
			ID:          "basic-auth",
			Name:        "Basic Authentication",
			Description: "HTTP basic authentication middleware",
			IsPremium:   false,
		},
		{
			ID:          "email",
			Name:        "Email Support",
			Description: "Helpers for sending emails via SMTP",
			IsPremium:   false,
		},
		{
			ID:          "error-notifications",
			Name:        "Error Notifications",
			Description: "Send error alerts to admin email",
			IsPremium:   false,
		},
		{
			ID:          "gitignore",
			Name:        "Gitignore",
			Description: "Common .gitignore file for Go projects",
			IsPremium:   false,
		},
		{
			ID:          "live-reload",
			Name:        "Live Reload",
			Description: "Auto-rebuild and restart during development",
			IsPremium:   false,
		},
		{
			ID:          "secure-cookies",
			Name:        "Secure Cookies",
			Description: "Signed and encrypted cookie support",
			IsPremium:   false,
		},
		{
			ID:          "sql-migrations",
			Name:        "SQL Migrations",
			Description: "Database migration tools",
			IsPremium:   false,
		},
		{
			ID:          "automatic-https",
			Name:        "Automatic HTTPS",
			Description: "TLS certificate management via Let's Encrypt",
			IsPremium:   true,
		},
		{
			ID:          "custom-error-pages",
			Name:        "Custom Error Pages",
			Description: "Custom HTML pages for error responses",
			IsPremium:   true,
		},
		{
			ID:          "user-accounts",
			Name:        "User Accounts",
			Description: "User authentication and management",
			IsPremium:   true,
		},
	}

	return features, nil
}

// Helper functions

// validateOptions validates the scaffold options
func (s *GeneratorServiceImpl) validateOptions(options model.ScaffoldOptions) error {
	if options.AppType != "api" && options.AppType != "webapp" {
		return errors.New("invalid application type")
	}

	if options.RouterType == "" {
		return errors.New("router type is required")
	}

	if options.ModulePath == "" {
		return errors.New("module path is required")
	}

	// Validate features
	for _, feature := range options.Features {
		// Check if the feature exists
		exists := false
		features, _ := s.GetAvailableFeatures()
		for _, f := range features {
			if f.ID == feature && !f.IsPremium {
				exists = true
				break
			}
		}

		if !exists {
			return fmt.Errorf("invalid feature: %s", feature)
		}
	}

	// Validate premium features
	for _, feature := range options.PremiumFeatures {
		// Check if the feature exists and is premium
		exists := false
		features, _ := s.GetAvailableFeatures()
		for _, f := range features {
			if f.ID == feature && f.IsPremium {
				exists = true
				break
			}
		}

		if !exists {
			return fmt.Errorf("invalid premium feature: %s", feature)
		}
	}

	return nil
}

// getTemplateForOptions returns the appropriate template for the provided options
func (s *GeneratorServiceImpl) getTemplateForOptions(options model.ScaffoldOptions) (*model.Template, error) {
	// This is a simplified implementation
	// In a real application, we would select the template based on the options
	templates, err := s.templateRepo.GetByType(options.AppType)
	if err != nil {
		return nil, err
	}

	if len(templates) == 0 {
		return nil, fmt.Errorf("no templates found for application type: %s", options.AppType)
	}

	// Find a template that matches the router type
	for _, tmpl := range templates {
		if tmpl.ID == options.AppType+"-"+options.RouterType {
			return tmpl, nil
		}
	}

	// If no exact match, return the first template
	return templates[0], nil
}

// processTemplate processes the template with the provided options
func (s *GeneratorServiceImpl) processTemplate(tmpl *model.Template, options model.ScaffoldOptions, outputDir string) error {
	// Create template data with all options and helper functions
	templateData := map[string]interface{}{
		"AppType":      options.AppType,
		"DatabaseType": options.DatabaseType,
		"RouterType":   options.RouterType,
		"ConfigType":   options.ConfigType,
		"LogFormat":    options.LogFormat,
		"ModulePath":   options.ModulePath,
		"Features":     options.Features,
		"Premium":      options.PremiumFeatures,
		// Helper functions
		"HasFeature": func(feature string) bool {
			for _, f := range options.Features {
				if f == feature {
					return true
				}
			}
			for _, f := range options.PremiumFeatures {
				if f == feature {
					return true
				}
			}
			return false
		},
	}

	// Walk through the template directory
	return filepath.Walk(tmpl.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if this is a template file
		isTemplate := filepath.Ext(path) == ".tmpl"

		// Get the relative path from the template root
		relPath, err := filepath.Rel(tmpl.Path, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// Determine the output file path
		outputPath := filepath.Join(outputDir, relPath)
		if isTemplate {
			// Remove .tmpl extension for template files
			outputPath = outputPath[:len(outputPath)-5]
		}

		// Create directory structure if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		// If this is a template file, process it
		if isTemplate {
			// Parse the template
			t, err := template.New(filepath.Base(path)).ParseFiles(path)
			if err != nil {
				return fmt.Errorf("failed to parse template %s: %w", path, err)
			}

			// Create output file
			out, err := os.Create(outputPath)
			if err != nil {
				return fmt.Errorf("failed to create output file %s: %w", outputPath, err)
			}
			defer out.Close()

			// Execute the template
			if err := t.Execute(out, templateData); err != nil {
				return fmt.Errorf("failed to execute template %s: %w", path, err)
			}
		} else {
			// For non-template files, just copy them
			srcFile, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open source file %s: %w", path, err)
			}
			defer srcFile.Close()

			dstFile, err := os.Create(outputPath)
			if err != nil {
				return fmt.Errorf("failed to create destination file %s: %w", outputPath, err)
			}
			defer dstFile.Close()

			if _, err := io.Copy(dstFile, srcFile); err != nil {
				return fmt.Errorf("failed to copy file %s: %w", path, err)
			}
		}

		return nil
	})
}

// createZipArchive creates a ZIP archive from the scaffold directory
func (s *GeneratorServiceImpl) createZipArchive(dir, zipPath string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()

	// Create a new zip writer
	zw := zip.NewWriter(zipFile)
	defer zw.Close()

	// Create a codebase/ directory in the zip root
	codebaseDir := "codebase/"
	codebaseHeader := &zip.FileHeader{
		Name:   codebaseDir,
		Method: zip.Deflate,
	}
	codebaseHeader.SetMode(0755 | os.ModeDir)
	_, err = zw.CreateHeader(codebaseHeader)
	if err != nil {
		return fmt.Errorf("failed to create codebase directory in zip: %w", err)
	}

	// Walk through the directory tree
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create a zip header based on the file info
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("failed to create file header: %w", err)
		}

		// Set the relative path as the name in the zip, but prefix with "codebase/"
		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// Skip the root directory
		if relPath == "." {
			return nil
		}

		// Add "codebase/" prefix
		if info.IsDir() {
			// Add trailing slash to folders
			header.Name = codebaseDir + relPath + "/"
		} else {
			header.Name = codebaseDir + relPath
		}

		// Set compression method
		header.Method = zip.Deflate

		if info.IsDir() {
			// For directories, just write the header
			_, err = zw.CreateHeader(header)
			return err
		}

		// Create the file in the zip
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("failed to create file in zip: %w", err)
		}

		// Open the file for reading
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file for reading: %w", err)
		}
		defer file.Close()

		// Copy the file contents to the zip
		_, err = io.Copy(writer, file)
		return err
	})

	if err != nil {
		return fmt.Errorf("failed to walk directory: %w", err)
	}

	return nil
}

// generateID generates a unique ID for a scaffold
func generateID() string {
	// Use crypto/rand for secure random generation
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	const length = 12

	b := make([]byte, length)
	for i := range b {
		// Get a random index within the charset
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}

	result := string(b)
	fmt.Printf("Generated random ID: %s\n", result) // Debug log

	return result
}

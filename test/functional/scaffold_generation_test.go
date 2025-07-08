package functional

import (
	"archive/zip"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestGenerateScaffoldAPI tests scaffold generation for API projects
func TestGenerateScaffoldAPI(t *testing.T) {
	// Skip test if running in CI environment
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping functional test in CI environment")
	}

	// Create temporary directory for test
	tmpDir, err := os.MkdirTemp("", "scaffold-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	outputZip := filepath.Join(tmpDir, "output.zip")
	
	// Run go-scaffold CLI
	cmd := exec.Command("go-scaffold", "generate",
		"--app-type", "api",
		"--router-type", "echo",
		"--database-type", "postgresql",
		"--config-type", "env",
		"--log-format", "json",
		"--module-path", "github.com/example/testapi",
		"--feature", "basic-auth",
		"--feature", "sql-migrations",
		"--output", outputZip)
	
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		// Try finding the binary in ../build directory
		binPath := filepath.Join("..", "..", "build", "go-scaffold")
		if _, err := os.Stat(binPath); os.IsNotExist(err) {
			t.Skipf("go-scaffold binary not found, skipping test: %v", err)
			return
		}
		
		cmd = exec.Command(binPath, "generate",
			"--app-type", "api",
			"--router-type", "echo",
			"--database-type", "postgresql",
			"--config-type", "env",
			"--log-format", "json",
			"--module-path", "github.com/example/testapi",
			"--feature", "basic-auth",
			"--feature", "sql-migrations",
			"--output", outputZip)
		
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to run go-scaffold: %v", err)
		}
	}
	
	// Check if zip file was created
	if _, err := os.Stat(outputZip); os.IsNotExist(err) {
		t.Fatalf("Output ZIP file was not created: %v", err)
	}
	
	// Extract ZIP file
	extractDir := filepath.Join(tmpDir, "extract")
	if err := os.Mkdir(extractDir, 0755); err != nil {
		t.Fatalf("Failed to create extraction directory: %v", err)
	}
	
	if err := extractZip(outputZip, extractDir); err != nil {
		t.Fatalf("Failed to extract ZIP file: %v", err)
	}
	
	// Check for required files
	requiredFiles := []string{
		"go.mod",
		"Makefile",
		"README.md",
		"cmd/api/main.go",
		"internal/config/config.go",
		"internal/handlers/api.go",
	}
	
	for _, file := range requiredFiles {
		path := findFile(extractDir, file)
		if path == "" {
			t.Errorf("Required file not found: %s", file)
			continue
		}
		
		// Check file content
		content, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("Failed to read file %s: %v", file, err)
			continue
		}
		
		contentStr := string(content)
		
		// Check for module path
		if strings.HasSuffix(file, "go.mod") && !strings.Contains(contentStr, "github.com/example/testapi") {
			t.Errorf("go.mod does not contain correct module path")
		}
		
		// Check for imports based on features
		if strings.HasSuffix(file, "config.go") && !strings.Contains(contentStr, "database") {
			t.Errorf("config.go does not contain database configuration")
		}
		
		// Check for PostgreSQL imports
		if strings.HasSuffix(file, "db.go") && !strings.Contains(contentStr, "postgres") {
			t.Errorf("db.go does not contain PostgreSQL imports")
		}
	}
	
	// Check if project builds
	buildDir := findCodebaseDir(extractDir)
	if buildDir == "" {
		t.Fatal("Could not find codebase directory")
	}
	
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = buildDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		t.Errorf("Failed to run go mod tidy: %v", err)
	}
	
	cmd = exec.Command("go", "build", "./...")
	cmd.Dir = buildDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		t.Errorf("Failed to build project: %v", err)
	}
}

// TestGenerateScaffoldWebApp tests scaffold generation for web app projects
func TestGenerateScaffoldWebApp(t *testing.T) {
	// Skip test if running in CI environment
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping functional test in CI environment")
	}

	// Create temporary directory for test
	tmpDir, err := os.MkdirTemp("", "scaffold-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	outputZip := filepath.Join(tmpDir, "output.zip")
	
	// Run go-scaffold CLI
	cmd := exec.Command("go-scaffold", "generate",
		"--app-type", "webapp",
		"--router-type", "chi",
		"--database-type", "sqlite",
		"--config-type", "flags",
		"--log-format", "json",
		"--module-path", "github.com/example/testwebapp",
		"--feature", "live-reload",
		"--output", outputZip)
	
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		// Try finding the binary in ../build directory
		binPath := filepath.Join("..", "..", "build", "go-scaffold")
		if _, err := os.Stat(binPath); os.IsNotExist(err) {
			t.Skipf("go-scaffold binary not found, skipping test: %v", err)
			return
		}
		
		cmd = exec.Command(binPath, "generate",
			"--app-type", "webapp",
			"--router-type", "chi",
			"--database-type", "sqlite",
			"--config-type", "flags",
			"--log-format", "json",
			"--module-path", "github.com/example/testwebapp",
			"--feature", "live-reload",
			"--output", outputZip)
		
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to run go-scaffold: %v", err)
		}
	}
	
	// Check if zip file was created
	if _, err := os.Stat(outputZip); os.IsNotExist(err) {
		t.Fatalf("Output ZIP file was not created: %v", err)
	}
	
	// Extract ZIP file
	extractDir := filepath.Join(tmpDir, "extract")
	if err := os.Mkdir(extractDir, 0755); err != nil {
		t.Fatalf("Failed to create extraction directory: %v", err)
	}
	
	if err := extractZip(outputZip, extractDir); err != nil {
		t.Fatalf("Failed to extract ZIP file: %v", err)
	}
	
	// Check for required files
	requiredFiles := []string{
		"go.mod",
		"Makefile",
		"README.md",
		"cmd/web/main.go",
		"internal/config/config.go",
		"internal/handlers/handlers.go",
	}
	
	for _, file := range requiredFiles {
		path := findFile(extractDir, file)
		if path == "" {
			t.Errorf("Required file not found: %s", file)
			continue
		}
		
		// Check file content
		content, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("Failed to read file %s: %v", file, err)
			continue
		}
		
		contentStr := string(content)
		
		// Check for module path
		if strings.HasSuffix(file, "go.mod") && !strings.Contains(contentStr, "github.com/example/testwebapp") {
			t.Errorf("go.mod does not contain correct module path")
		}
		
		// Check for imports based on database type
		if strings.HasSuffix(file, "config.go") && !strings.Contains(contentStr, "sqlite") {
			t.Errorf("config.go does not contain SQLite configuration")
		}
	}
	
	// Check if project builds
	buildDir := findCodebaseDir(extractDir)
	if buildDir == "" {
		t.Fatal("Could not find codebase directory")
	}
	
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = buildDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		t.Errorf("Failed to run go mod tidy: %v", err)
	}
	
	cmd = exec.Command("go", "build", "./...")
	cmd.Dir = buildDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		t.Errorf("Failed to build project: %v", err)
	}
}

// Helper functions

// extractZip extracts a ZIP file to a destination directory
func extractZip(zipFile, destDir string) error {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer r.Close()
	
	for _, f := range r.File {
		fPath := filepath.Join(destDir, f.Name)
		
		// Create directory if needed
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fPath, os.ModePerm); err != nil {
				return err
			}
			continue
		}
		
		// Create parent directory if needed
		if err := os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
			return err
		}
		
		// Create file
		outFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		
		// Extract file content
		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}
		
		if _, err := io.Copy(outFile, rc); err != nil {
			outFile.Close()
			rc.Close()
			return err
		}
		
		outFile.Close()
		rc.Close()
	}
	
	return nil
}

// findFile looks for a file in a directory and its subdirectories
func findFile(root, fileName string) string {
	var result string
	
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		
		if !info.IsDir() && filepath.Base(path) == fileName {
			result = path
			return filepath.SkipDir
		}
		
		return nil
	})
	
	return result
}

// findCodebaseDir finds the directory containing the go.mod file
func findCodebaseDir(root string) string {
	var result string
	
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		
		if !info.IsDir() && filepath.Base(path) == "go.mod" {
			result = filepath.Dir(path)
			return filepath.SkipDir
		}
		
		return nil
	})
	
	return result
}

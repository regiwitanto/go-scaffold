package testutil

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// FindProjectRoot returns the absolute path to the project root
func FindProjectRoot() (string, error) {
	// Start from the current file's directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil
	}

	// Navigate up to find the project root (where go.mod is)
	dir := filepath.Dir(filename)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached the filesystem root
			return "", nil
		}
		dir = parent
	}
}

// AssertFileExists checks if a file exists
func AssertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("File does not exist: %s", path)
	}
}

// AssertFileContains checks if a file contains a string
func AssertFileContains(t *testing.T, path, content string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Errorf("Failed to read file %s: %v", path, err)
		return
	}

	if !strings.Contains(string(data), content) {
		t.Errorf("File %s does not contain expected content: %s", path, content)
	}
}

// AssertDirExists checks if a directory exists
func AssertDirExists(t *testing.T, path string) {
	t.Helper()
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		t.Errorf("Directory does not exist: %s", path)
		return
	}

	if !info.IsDir() {
		t.Errorf("Path exists but is not a directory: %s", path)
	}
}

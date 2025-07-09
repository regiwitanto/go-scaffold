# Go-Scaffold Testing Guide

This document provides detailed information about the Go-Scaffold testing strategy, including instructions for running tests, creating new tests, and best practices for maintaining and extending the test suite.

## Overview of Testing Strategy

The testing strategy for Go-Scaffold follows the testing pyramid approach, with:

1. **Unit tests** - Testing individual components in isolation
2. **Integration tests** - Testing interactions between components
3. **Functional tests** - Testing the entire system from an external perspective

## Test Directory Structure

```
test/
├── mocks/                   # Mock objects for testing
│   ├── mock_generator_service.go
│   ├── mock_scaffold_repository.go
│   └── mock_template_repository.go
├── testutil/                # Common test utilities
│   ├── helper.go
│   └── template_util.go
├── unit/                    # Unit tests
│   ├── api_handlers_test.go
│   ├── api_handlers_extended_test.go
│   ├── benchmark_test.go
│   └── domain_model_test.go
├── integration/             # Integration tests
│   ├── template_consistency_test.go
│   └── template_content_test.go
├── functional/              # Functional tests
│   └── scaffold_generation_test.go
├── run-tests.sh             # Test runner script
├── README.md                # Test documentation
└── main_test.go             # Test entry point
```

## Running Tests

You can run tests using several methods:

### 1. Using Make targets

```bash
make test             # Run all tests
make test-unit        # Run unit tests only
make test-integration # Run integration tests only
make test-functional  # Run functional tests only
make test-cover       # Run tests with coverage
```

### 2. Using the run-tests.sh script

```bash
./test/run-tests.sh --all          # Run all tests
./test/run-tests.sh --unit         # Run unit tests only
./test/run-tests.sh --integration  # Run integration tests only
./test/run-tests.sh --functional   # Run functional tests only
./test/run-tests.sh --benchmark    # Run benchmark tests
./test/run-tests.sh --cover        # Generate coverage report
```

### 3. Using Go test directly

```bash
go test -v ./test/unit/...           # Run unit tests
go test -v ./test/integration/...    # Run integration tests
go test -v ./test/functional/...     # Run functional tests
go test -bench=. ./test/unit/...     # Run benchmark tests
```

## Writing New Tests

### Unit Tests

Unit tests should focus on testing individual components in isolation. They should be fast, deterministic, and not rely on external dependencies.

Example of a unit test:

```go
func TestScaffoldOptions(t *testing.T) {
    // Setup test data
    opts := model.ScaffoldOptions{
        AppType:      "api",
        DatabaseType: "postgresql",
        RouterType:   "echo",
        ConfigType:   "env",
        LogFormat:    "json",
        ModulePath:   "github.com/testuser/testproject",
        Features:     []string{"migrations", "logging"},
    }

    // Test behavior
    jsonData, err := json.Marshal(opts)
    if err != nil {
        t.Fatalf("Failed to marshal ScaffoldOptions: %v", err)
    }

    // Validate results
    var unmarshalledOpts model.ScaffoldOptions
    err = json.Unmarshal(jsonData, &unmarshalledOpts)
    if err != nil {
        t.Fatalf("Failed to unmarshal ScaffoldOptions: %v", err)
    }

    // Compare fields
    if opts.AppType != unmarshalledOpts.AppType {
        t.Errorf("AppType does not match: expected %s, got %s", 
                 opts.AppType, unmarshalledOpts.AppType)
    }
}
```

### Integration Tests

Integration tests should focus on testing interactions between components. They may rely on external dependencies such as the filesystem, but should be designed to minimize external dependencies.

Example of an integration test:

```go
func TestTemplateRendering(t *testing.T) {
    // Find project root
    rootDir, err := testutil.FindProjectRoot()
    if err != nil {
        t.Fatalf("Failed to find project root: %v", err)
    }

    // Construct absolute template path
    tmplPath := filepath.Join(rootDir, "templates/api/echo/go.mod.tmpl")
    
    // Parse template
    tmpl, err := template.ParseFiles(tmplPath)
    if err != nil {
        t.Fatalf("Failed to parse template: %v", err)
    }
    
    // Execute template with test data
    output := new(strings.Builder)
    err = tmpl.Execute(output, testutil.NewTemplateData())
    if err != nil {
        t.Fatalf("Failed to execute template: %v", err)
    }
    
    // Validate output
    if !strings.Contains(output.String(), "github.com/example/testapp") {
        t.Errorf("Template output does not contain expected module path")
    }
}
```

### Functional Tests

Functional tests should focus on testing the entire system from an external perspective. They may rely on external dependencies such as the filesystem and command execution, but should be designed to be as self-contained as possible.

Example of a functional test:

```go
func TestGenerateScaffoldAPI(t *testing.T) {
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
        "--module-path", "github.com/example/testapi",
        "--feature", "basic-auth",
        "--output", outputZip)
    
    if err := cmd.Run(); err != nil {
        t.Fatalf("Failed to run go-scaffold: %v", err)
    }
    
    // Verify output
    if _, err := os.Stat(outputZip); os.IsNotExist(err) {
        t.Fatalf("Output ZIP file was not created: %v", err)
    }
}
```

## Test Mocks

The `test/mocks` directory contains mock implementations of interfaces used in the application. These mocks can be used in tests to simulate the behavior of dependencies.

Example of using a mock:

```go
// Create mock repositories
mockTemplateRepo := &mocks.MockTemplateRepository{}
mockScaffoldRepo := &mocks.MockScaffoldRepository{}

// Setup specific behavior
mockTemplateRepo.GetByTypeFunc = func(templateType string) ([]*model.Template, error) {
    if templateType == "api" {
        return []*model.Template{
            {
                ID:          "api-echo",
                Name:        "API with Echo",
                Path:        "templates/api/echo",
                Type:        "api",
            },
        }, nil
    }
    return nil, errors.New("template not found")
}

// Create service with mocks
generatorService := service.NewGeneratorService(mockTemplateRepo, mockScaffoldRepo, "/tmp")
```

## Test Utilities

The `test/testutil` directory contains common utilities for testing. These utilities can be used to simplify test code and avoid duplication.

Example of using test utilities:

```go
// Find project root
rootDir, err := testutil.FindProjectRoot()
if err != nil {
    t.Fatalf("Failed to find project root: %v", err)
}

// Create template data
data := testutil.NewTemplateData().WithDatabaseType("postgresql")

// Render template
content, err := testutil.RenderTemplate(t, templatePath, data)
if err != nil {
    t.Fatalf("Failed to render template: %v", err)
}
```

## Best Practices

1. **Test one thing at a time**: Each test should focus on testing one specific behavior or aspect of the system.

2. **Use descriptive test names**: Test names should clearly describe what the test is testing.

3. **Use helper functions**: Extract common setup and teardown code into helper functions.

4. **Clean up after tests**: Use `defer` to ensure cleanup happens even if the test fails.

5. **Use table-driven tests**: Table-driven tests can help test multiple cases with minimal code duplication.

6. **Test edge cases**: Make sure to test edge cases and error conditions, not just the happy path.

7. **Keep tests fast**: Tests should run quickly to encourage frequent testing.

8. **Keep tests independent**: Tests should not depend on the results of other tests.

9. **Use assertions sparingly**: Use assertions only for the specific behavior being tested, not for side effects.

10. **Keep tests simple**: Complex tests are hard to maintain and can become fragile.

## Continuous Integration

Tests are automatically run as part of the continuous integration (CI) process. The CI pipeline includes running all tests, generating a coverage report, and failing the build if any tests fail.

The CI pipeline also includes running the tests with the race detector enabled to catch race conditions that might not be detected during normal testing.

## Current Test Status

### Integration Tests

The integration tests in `/test/integration` are currently skipped as they need further work to align with the actual structure of the template files in the codebase. These tests validate that templates are consistent, render correctly, and follow best practices.

#### Issues to Fix

1. **Template Path Structure**:
   - Tests expect `main.go.tmpl` files at the root of the template directories, but they are actually in subdirectories like `cmd/api/main.go.tmpl` or `cmd/web/main.go.tmpl`
   - Some templates are missing entirely in certain router types (gin/standard webapp templates)

2. **Template Function Support**:
   - Templates use `.HasFeature` function calls, which needs to be properly implemented
   - Need to add function map to template rendering with `HasFeature` implementation

3. **Missing Template Data Fields**:
   - Templates reference fields like `.Binary`, `.Subject`, `.Timestamp` that don't exist in the test TemplateData struct
   - Added fields to TemplateData but still need proper integration with template function calls

#### Plans to Re-enable

To re-enable integration tests:

1. Update `testutil.TemplateData` to include all required fields used in templates
2. Modify template path checks to match the actual directory structure
3. Ensure template function map includes all functions used in templates:
   - `HasFeature`
   - Any other custom functions used in templates
4. Update test case assertions to match actual template outputs

See the detailed README in the `/test/integration` directory for more information.

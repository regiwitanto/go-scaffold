// build.go - A Go-based build system to replace the Makefile
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	appName  = "go-scaffold"
	version  = "1.0.0"
	buildDir = "build"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		help()
		return
	}

	switch args[0] {
	case "all":
		clean()
		build()
	case "build":
		build()
	case "clean":
		clean()
	case "test":
		test("")
	case "test-unit":
		test("./test/unit/...")
	case "test-integration":
		test("./test/integration/...")
	case "test-functional":
		test("./test/functional/...")
	case "test-cover":
		testCover()
	case "run":
		run(false)
	case "run-dev":
		run(true)
	case "api-docs":
		apiDocs()
	case "help":
		help()
	default:
		fmt.Printf("Unknown command: %s\n", args[0])
		help()
	}
}

func execute(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func build() {
	fmt.Printf("Building %s...\n", appName)
	err := os.MkdirAll(buildDir, 0755)
	if err != nil {
		fmt.Printf("Error creating build directory: %v\n", err)
		os.Exit(1)
	}

	err = execute("go", "build", "-o", filepath.Join(buildDir, appName), "-v", "./cmd/scaffold")
	if err != nil {
		fmt.Printf("Error building: %v\n", err)
		os.Exit(1)
	}
}

func clean() {
	fmt.Println("Cleaning...")
	err := os.RemoveAll(buildDir)
	if err != nil {
		fmt.Printf("Error cleaning: %v\n", err)
		os.Exit(1)
	}
}

func test(path string) {
	if path == "" {
		fmt.Println("Running all tests...")
		path = "./..."
	} else {
		fmt.Printf("Running tests in %s...\n", path)
	}

	err := execute("go", "test", "-v", path)
	if err != nil {
		fmt.Printf("Tests failed: %v\n", err)
		os.Exit(1)
	}
}

func testCover() {
	fmt.Println("Running tests with coverage...")
	err := execute("go", "test", "-cover", "-v", "./...")
	if err != nil {
		fmt.Printf("Tests failed: %v\n", err)
		os.Exit(1)
	}

	err = execute("go", "test", "-coverprofile=coverage.out", "./...")
	if err != nil {
		fmt.Printf("Coverage generation failed: %v\n", err)
		os.Exit(1)
	}

	err = execute("go", "tool", "cover", "-func=coverage.out")
	if err != nil {
		fmt.Printf("Coverage report failed: %v\n", err)
		os.Exit(1)
	}
}

func run(dev bool) {
	if dev {
		fmt.Printf("Running %s in development mode...\n", appName)
		// Check if air is installed
		_, err := exec.LookPath("air")
		if err != nil {
			fmt.Println("Error: 'air' is not installed or not in PATH")
			fmt.Println("Install it with: go install github.com/cosmtrek/air@latest")
			os.Exit(1)
		}
		err = execute("air")
		if err != nil {
			fmt.Printf("Error running with air: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Running %s...\n", appName)
		err := execute("go", "run", "./cmd/scaffold")
		if err != nil {
			fmt.Printf("Error running: %v\n", err)
			os.Exit(1)
		}
	}
}

func apiDocs() {
	fmt.Println("API documentation is available at http://localhost:8081/api-docs")
	fmt.Println("OpenAPI spec is available at http://localhost:8081/api/docs")
}

func help() {
	fmt.Println("Usage: go run build.go [command]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  all             Clean and build the application")
	fmt.Println("  build           Build the application")
	fmt.Println("  clean           Clean build artifacts")
	fmt.Println("  test            Run all tests")
	fmt.Println("  test-unit       Run unit tests only")
	fmt.Println("  test-integration Run integration tests only")
	fmt.Println("  test-functional Run functional tests only")
	fmt.Println("  test-cover      Run tests with coverage")
	fmt.Println("  run             Run the application")
	fmt.Println("  run-dev         Run the application in development mode (requires air)")
	fmt.Println("  api-docs        Show API documentation information")
	fmt.Println("  help            Show this help")
}

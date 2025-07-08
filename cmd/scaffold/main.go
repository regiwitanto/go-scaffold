package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/regiwitanto/go-scaffold/internal/application/service"
	"github.com/regiwitanto/go-scaffold/internal/infrastructure/storage/scaffold"
	"github.com/regiwitanto/go-scaffold/internal/infrastructure/storage/template"
	"github.com/regiwitanto/go-scaffold/internal/interfaces/api/handler"
	"github.com/regiwitanto/go-scaffold/internal/interfaces/api/routes"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it. Using environment variables.")
	}

	// Initialize Echo instance
	e := echo.New()
	e.HideBanner = true

	// Create assets directory if it doesn't exist
	if err := os.MkdirAll("assets", 0755); err != nil {
		log.Fatalf("Failed to create assets directory: %v", err)
	}

	// Create assets/js directory if it doesn't exist
	if err := os.MkdirAll("assets/js", 0755); err != nil {
		log.Fatalf("Failed to create assets/js directory: %v", err)
	}

	// Set up temp directory for scaffold generation
	tempDir := os.Getenv("TEMP_DIR")
	if tempDir == "" {
		tempDir = filepath.Join(os.TempDir(), "go-scaffold")
	}
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		log.Fatalf("Failed to create temp directory: %v", err)
	}

	// Create the templates directory structure if it doesn't exist
	templatesDir := os.Getenv("TEMPLATE_DIR")
	if templatesDir == "" {
		templatesDir = "templates"
	}
	templateDirs := []string{
		filepath.Join(templatesDir, "api", "chi"),
		filepath.Join(templatesDir, "api", "echo"),
		filepath.Join(templatesDir, "api", "gin"),
		filepath.Join(templatesDir, "api", "standard"),
		filepath.Join(templatesDir, "webapp", "chi"),
		filepath.Join(templatesDir, "webapp", "echo"),
		filepath.Join(templatesDir, "webapp", "gin"),
		filepath.Join(templatesDir, "webapp", "standard"),
		filepath.Join(templatesDir, "web"),
	}

	for _, dir := range templateDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Failed to create template directory %s: %v", dir, err)
		}
	}

	// Initialize repositories
	templateRepo, err := template.NewFilesystemRepository("templates")
	if err != nil {
		log.Fatalf("Failed to initialize template repository: %v", err)
	}
	scaffoldRepo := scaffold.NewInMemoryRepository()

	// Initialize services
	generatorService := service.NewGeneratorService(templateRepo, scaffoldRepo, tempDir)

	// Initialize handlers
	generatorHandler := handler.NewGeneratorHandler(generatorService)

	// Setup routes
	routes.SetupRoutes(e, generatorHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Clear message to show where the server is running
	serverURL := fmt.Sprintf("http://localhost:%s", port)
	log.Printf("Starting server on %s", serverURL)
	log.Printf("Web UI available at %s", serverURL)
	log.Printf("API Documentation at %s/api-docs", serverURL)

	e.Logger.Fatal(e.Start(":" + port))
}

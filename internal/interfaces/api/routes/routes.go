package routes

import (
	"net/http"

	"github.com/regiwitanto/go-scaffold/internal/interfaces/api/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(e *echo.Echo, generatorHandler *handler.GeneratorHandler) {
	// Basic middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Redirect root to API docs
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/api/docs")
	})

	// API routes
	api := e.Group("/api")
	{
		api.GET("/health", generatorHandler.HandleHealthCheck)
		api.GET("/features", generatorHandler.HandleListFeatures)
		api.GET("/templates", generatorHandler.HandleListTemplates)
		api.POST("/generate", generatorHandler.HandleGenerateScaffold)
		api.GET("/download/:id", generatorHandler.HandleDownloadScaffold)

		// API documentation
		apiDocsHandler := handler.NewApiDocsHandler()
		api.GET("/docs", apiDocsHandler.HandleApiDocs)
	}

	// Swagger UI endpoint - replaced with JSON response since we're removing HTML
	e.GET("/api-docs", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message":  "API documentation available at /api/docs",
			"docs_url": "/api/docs",
		})
	})

	// 404 handler
	e.Any("/*", func(c echo.Context) error {
		return c.HTML(http.StatusNotFound, `
			<html>
				<head>
					<title>Page Not Found</title>
					<style>
						body { font-family: Arial, sans-serif; line-height: 1.6; margin: 0; padding: 20px; color: #333; }
						.container { max-width: 800px; margin: 0 auto; text-align: center; }
						h1 { color: #0077cc; }
						.error { font-size: 120px; margin: 0; }
					</style>
				</head>
				<body>
					<div class="container">
						<h1 class="error">404</h1>
						<h2>Page Not Found</h2>
						<p>The page you are looking for does not exist.</p>
						<p><a href="/">Go back to home page</a></p>
					</div>
				</body>
			</html>
		`)
	})
}

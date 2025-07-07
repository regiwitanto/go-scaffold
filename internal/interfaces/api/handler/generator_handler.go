package handler

import (
	"net/http"

	"github.com/regiwitanto/echo-scaffold/internal/domain/model"
	"github.com/regiwitanto/echo-scaffold/internal/domain/service"

	"github.com/labstack/echo/v4"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Status string `json:"status" example:"OK"`
}

// GenerateResponse represents a successful scaffold generation response
type GenerateResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

// GeneratorHandler handles API requests related to scaffold generation
type GeneratorHandler struct {
	generatorService service.GeneratorService
}

// NewGeneratorHandler creates a new generator handler
func NewGeneratorHandler(generatorService service.GeneratorService) *GeneratorHandler {
	return &GeneratorHandler{
		generatorService: generatorService,
	}
}

// HandleHealthCheck godoc
// @Summary Health check endpoint
// @Description Check if the API is running
// @Tags system
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *GeneratorHandler) HandleHealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, HealthResponse{
		Status: "OK",
	})
}

// HandleListFeatures godoc
// @Summary List features endpoint
// @Description Get a list of all available features
// @Tags features
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} ErrorResponse
// @Router /features [get]
func (h *GeneratorHandler) HandleListFeatures(c echo.Context) error {
	features, err := h.generatorService.GetAvailableFeatures()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to retrieve features",
		})
	}

	// Separate regular and premium features
	regularFeatures := make([]*model.Feature, 0)
	premiumFeatures := make([]*model.Feature, 0)

	for _, feature := range features {
		if feature.IsPremium {
			premiumFeatures = append(premiumFeatures, feature)
		} else {
			regularFeatures = append(regularFeatures, feature)
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"features":        regularFeatures,
		"premiumFeatures": premiumFeatures,
	})
}

// HandleListTemplates godoc
// @Summary List templates endpoint
// @Description Get a list of all available templates
// @Tags templates
// @Produce json
// @Success 200 {array} model.Template
// @Failure 500 {object} ErrorResponse
// @Router /templates [get]
func (h *GeneratorHandler) HandleListTemplates(c echo.Context) error {
	templates, err := h.generatorService.GetAllTemplates()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to retrieve templates",
		})
	}

	return c.JSON(http.StatusOK, templates)
}

// HandleGenerateScaffold godoc
// @Summary Generate scaffold endpoint
// @Description Generate a scaffold project based on provided options
// @Tags generator
// @Accept json
// @Produce json
// @Param options body model.ScaffoldOptions true "Scaffold options"
// @Success 200 {object} GenerateResponse
// @Failure 400 {object} ErrorResponse
// @Router /generate [post]
func (h *GeneratorHandler) HandleGenerateScaffold(c echo.Context) error {
	// Parse request body
	options := new(model.ScaffoldOptions)
	if err := c.Bind(options); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
		})
	}

	// Generate scaffold
	scaffold, err := h.generatorService.GenerateScaffold(*options)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, GenerateResponse{
		ID:      scaffold.ID,
		Message: "Scaffold generated successfully",
	})
}

// HandleDownloadScaffold godoc
// @Summary Download scaffold endpoint
// @Description Download a generated scaffold project by ID
// @Tags generator
// @Produce octet-stream
// @Param id path string true "Scaffold ID"
// @Success 200 {file} file
// @Failure 404 {object} ErrorResponse
// @Router /download/{id} [get]
func (h *GeneratorHandler) HandleDownloadScaffold(c echo.Context) error {
	id := c.Param("id")

	scaffold, err := h.generatorService.GetScaffold(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "Scaffold not found",
		})
	}

	return c.File(scaffold.FilePath)
}

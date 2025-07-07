package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// ApiDocsHandler handles API documentation
type ApiDocsHandler struct{}

// NewApiDocsHandler creates a new API docs handler
func NewApiDocsHandler() *ApiDocsHandler {
	return &ApiDocsHandler{}
}

// HandleApiDocs returns the OpenAPI specification for the API
func (h *ApiDocsHandler) HandleApiDocs(c echo.Context) error {
	docs := map[string]interface{}{
		"openapi": "3.0.0",
		"info": map[string]interface{}{
			"title":       "Go Scaffold Generator API",
			"description": "API for generating Go application scaffolds",
			"version":     "1.0.0",
		},
		"servers": []map[string]string{
			{
				"url": "/api",
			},
		},
		"paths": map[string]interface{}{
			"/health": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":     "Health Check",
					"description": "Returns the health status of the API",
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "OK",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"status": map[string]string{
												"type":    "string",
												"example": "OK",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"/templates": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":     "List Templates",
					"description": "Returns a list of available templates",
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "OK",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type":  "array",
										"items": templateSchema,
									},
								},
							},
						},
					},
				},
			},
			"/features": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":     "List Features",
					"description": "Returns a list of available features",
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "OK",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"features": map[string]interface{}{
												"type":  "array",
												"items": featureSchema,
											},
											"premiumFeatures": map[string]interface{}{
												"type":  "array",
												"items": featureSchema,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"/generate": map[string]interface{}{
				"post": map[string]interface{}{
					"summary":     "Generate Scaffold",
					"description": "Generates a scaffold based on provided options",
					"requestBody": map[string]interface{}{
						"required": true,
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": scaffoldOptionsSchema,
							},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "OK",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"id": map[string]string{
												"type":    "string",
												"example": "123",
											},
											"message": map[string]string{
												"type":    "string",
												"example": "Scaffold generated successfully",
											},
										},
									},
								},
							},
						},
						"400": map[string]interface{}{
							"description": "Bad Request",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"error": map[string]string{
												"type":    "string",
												"example": "Invalid request",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"/download/{id}": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":     "Download Scaffold",
					"description": "Downloads a generated scaffold by ID",
					"parameters": []map[string]interface{}{
						{
							"in":          "path",
							"name":        "id",
							"required":    true,
							"description": "Scaffold ID",
							"schema": map[string]string{
								"type": "string",
							},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "OK",
							"content": map[string]interface{}{
								"application/zip": map[string]interface{}{
									"schema": map[string]interface{}{
										"type":   "string",
										"format": "binary",
									},
								},
							},
						},
						"404": map[string]interface{}{
							"description": "Not Found",
							"content": map[string]interface{}{
								"application/json": map[string]interface{}{
									"schema": map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"error": map[string]string{
												"type":    "string",
												"example": "Scaffold not found",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return c.JSON(http.StatusOK, docs)
}

var templateSchema = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"id": map[string]string{
			"type":    "string",
			"example": "api-echo",
		},
		"name": map[string]string{
			"type":    "string",
			"example": "API with Echo router",
		},
		"description": map[string]string{
			"type":    "string",
			"example": "A RESTful API using the Echo framework",
		},
		"type": map[string]string{
			"type":    "string",
			"example": "api",
		},
		"path": map[string]string{
			"type":    "string",
			"example": "/templates/api/echo",
		},
	},
}

var featureSchema = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"id": map[string]string{
			"type":    "string",
			"example": "auth",
		},
		"name": map[string]string{
			"type":    "string",
			"example": "Authentication",
		},
		"description": map[string]string{
			"type":    "string",
			"example": "Basic authentication middleware",
		},
		"isPremium": map[string]interface{}{
			"type":    "boolean",
			"example": false,
		},
	},
}

var scaffoldOptionsSchema = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"appType": map[string]string{
			"type":    "string",
			"example": "api",
		},
		"routerType": map[string]string{
			"type":    "string",
			"example": "echo",
		},
		"databaseType": map[string]string{
			"type":    "string",
			"example": "postgresql",
		},
		"configType": map[string]string{
			"type":    "string",
			"example": "env",
		},
		"logFormat": map[string]string{
			"type":    "string",
			"example": "json",
		},
		"modulePath": map[string]string{
			"type":    "string",
			"example": "github.com/username/project",
		},
		"features": map[string]interface{}{
			"type": "array",
			"items": map[string]interface{}{
				"type":    "string",
				"example": "auth",
			},
		},
		"premiumFeatures": map[string]interface{}{
			"type": "array",
			"items": map[string]interface{}{
				"type":    "string",
				"example": "automatic-https",
			},
		},
	},
}

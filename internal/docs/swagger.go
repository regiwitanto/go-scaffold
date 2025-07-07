// Package docs generates Swagger documentation for the API
package docs

import "github.com/swaggo/swag"

// Registration of swagger docs
func init() {
	swag.Register(swag.Name, &swag.Spec{
		InfoInstanceName: "swagger",
		SwaggerTemplate:  docTemplate,
	})
}

const docTemplate = `{
    "swagger": "2.0",
    "info": {
        "description": "API service for Go Scaffold Generator.",
        "title": "Go Scaffold Generator API",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "email": "support@example.com"
        },
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        },
        "version": "1.0"
    },
    "host": "localhost:8081",
    "basePath": "/api",
    "paths": {
        "/health": {
            "get": {
                "description": "Check if the API is running",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "system"
                ],
                "summary": "Health check endpoint",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.HealthResponse"
                        }
                    }
                }
            }
        },
        "/templates": {
            "get": {
                "description": "Get a list of all available templates",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "templates"
                ],
                "summary": "List templates endpoint",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Template"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/features": {
            "get": {
                "description": "Get a list of all available features",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "features"
                ],
                "summary": "List features endpoint",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "properties": {
                                "features": {
                                    "type": "array",
                                    "items": {
                                        "$ref": "#/definitions/model.Feature"
                                    }
                                },
                                "premiumFeatures": {
                                    "type": "array",
                                    "items": {
                                        "$ref": "#/definitions/model.Feature"
                                    }
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/generate": {
            "post": {
                "description": "Generate a scaffold project based on provided options",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "generator"
                ],
                "summary": "Generate scaffold endpoint",
                "parameters": [
                    {
                        "description": "Scaffold options",
                        "name": "options",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.ScaffoldOptions"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.GenerateResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/download/{id}": {
            "get": {
                "description": "Download a generated scaffold project by ID",
                "produces": [
                    "application/zip"
                ],
                "tags": [
                    "generator"
                ],
                "summary": "Download scaffold endpoint",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Scaffold ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "file"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.ScaffoldOptions": {
            "type": "object",
            "properties": {
                "appType": {
                    "type": "string",
                    "example": "api"
                },
                "configType": {
                    "type": "string",
                    "example": "env"
                },
                "databaseType": {
                    "type": "string",
                    "example": "postgresql"
                },
                "features": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": ["auth", "email"]
                },
                "logFormat": {
                    "type": "string",
                    "example": "json"
                },
                "modulePath": {
                    "type": "string",
                    "example": "github.com/username/project"
                },
                "premiumFeatures": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "routerType": {
                    "type": "string",
                    "example": "echo"
                }
            }
        },
        "model.Template": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "path": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "model.Feature": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "isPremium": {
                    "type": "boolean"
                }
            }
        },
        "handler.HealthResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string",
                    "example": "OK"
                }
            }
        },
        "handler.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "handler.GenerateResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

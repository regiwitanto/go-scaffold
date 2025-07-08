// Package api provides API documentation and schema for the Go Scaffold Generator
package api

import (
	"github.com/regiwitanto/go-scaffold/internal/domain/model"
)

// HealthResponse represents a health check response
type HealthResponse struct {
	Status string `json:"status" example:"OK"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// GenerateResponse represents a successful scaffold generation response
type GenerateResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

// ScaffoldRequest represents a request to generate a scaffold
type ScaffoldRequest struct {
	Options model.ScaffoldOptions `json:"options"`
}

// FeatureResponse represents a response with regular and premium features
type FeatureResponse struct {
	Features        []*model.Feature `json:"features"`
	PremiumFeatures []*model.Feature `json:"premiumFeatures"`
}

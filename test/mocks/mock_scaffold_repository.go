package mocks

import (
	"github.com/regiwitanto/go-scaffold/internal/domain/model"
)

// MockScaffoldRepository is a mock implementation of the ScaffoldRepository interface
type MockScaffoldRepository struct {
	// Mock behavior functions
	SaveFunc    func(scaffold *model.GeneratedScaffold) error
	GetByIDFunc func(id string) (*model.GeneratedScaffold, error)
	DeleteFunc  func(id string) error

	// Tracking calls
	SaveCalled    bool
	SaveArg       *model.GeneratedScaffold
	GetByIDCalled bool
	GetByIDArg    string
	DeleteCalled  bool
	DeleteArg     string
}

// Save implements the ScaffoldRepository interface
func (m *MockScaffoldRepository) Save(scaffold *model.GeneratedScaffold) error {
	m.SaveCalled = true
	m.SaveArg = scaffold
	if m.SaveFunc != nil {
		return m.SaveFunc(scaffold)
	}
	return nil
}

// GetByID implements the ScaffoldRepository interface
func (m *MockScaffoldRepository) GetByID(id string) (*model.GeneratedScaffold, error) {
	m.GetByIDCalled = true
	m.GetByIDArg = id
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}

	// Return a mock scaffold
	return &model.GeneratedScaffold{
		ID: id,
		Options: model.ScaffoldOptions{
			AppType:      "api",
			RouterType:   "echo",
			DatabaseType: "postgresql",
			ModulePath:   "github.com/example/test-app",
		},
		CreatedAt: "2023-01-01T00:00:00Z",
		FilePath:  "/tmp/scaffolds/" + id + ".zip",
		Size:      1000,
	}, nil
}

// Delete implements the ScaffoldRepository interface
func (m *MockScaffoldRepository) Delete(id string) error {
	m.DeleteCalled = true
	m.DeleteArg = id
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

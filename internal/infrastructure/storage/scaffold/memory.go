package scaffold

import (
	"errors"
	"sync"

	"github.com/regiwitanto/echo-scaffold/internal/domain/model"
	"github.com/regiwitanto/echo-scaffold/internal/domain/repository"
)

// InMemoryRepository implements the ScaffoldRepository interface
// using in-memory storage
type InMemoryRepository struct {
	scaffolds map[string]*model.GeneratedScaffold
	mutex     sync.RWMutex
}

// NewInMemoryRepository creates a new in-memory scaffold repository
func NewInMemoryRepository() repository.ScaffoldRepository {
	return &InMemoryRepository{
		scaffolds: make(map[string]*model.GeneratedScaffold),
	}
}

// Save stores a generated scaffold
func (r *InMemoryRepository) Save(scaffold *model.GeneratedScaffold) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.scaffolds[scaffold.ID] = scaffold
	return nil
}

// GetByID returns a generated scaffold by ID
func (r *InMemoryRepository) GetByID(id string) (*model.GeneratedScaffold, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	scaffold, ok := r.scaffolds[id]
	if !ok {
		return nil, errors.New("scaffold not found")
	}

	return scaffold, nil
}

// Delete removes a generated scaffold
func (r *InMemoryRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.scaffolds[id]; !ok {
		return errors.New("scaffold not found")
	}

	delete(r.scaffolds, id)
	return nil
}

package repository

import (
	"context"
	"errors"
	"sync"

	"job-tracker-api/internal/model"
)

var (
	ErrApplicationNotFound = errors.New("application not found")
)

type MemoryApplicationRepository struct {
	mu           sync.RWMutex
	applications map[int64]model.Application
	nextID       int64
}

func NewMemoryApplicationRepository() *MemoryApplicationRepository {
	return &MemoryApplicationRepository{
		applications: make(map[int64]model.Application),
		nextID:       1,
	}
}

func (repository *MemoryApplicationRepository) Create(
	_ context.Context,
	application model.Application,
) (model.Application, error) {

	repository.mu.Lock()
	defer repository.mu.Unlock()

	application.ID = repository.nextID
	repository.nextID++

	repository.applications[application.ID] = application

	return application, nil
}

func (repository *MemoryApplicationRepository) GetByID(
	_ context.Context,
	id int64,
) (model.Application, error) {

	repository.mu.RLock()
	defer repository.mu.RUnlock()

	app, ok := repository.applications[id]
	if !ok {
		return model.Application{}, ErrApplicationNotFound
	}

	return app, nil

}

func (repository *MemoryApplicationRepository) Delete(
	_ context.Context,
	id int64,
) error {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	if _, ok := repository.applications[id]; !ok {
		return ErrApplicationNotFound
	}
	delete(repository.applications, id)
	return nil
}

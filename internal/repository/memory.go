package repository

import (
	"context"
	"errors"
	"sort"
	"strings"
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

func (repository *MemoryApplicationRepository) List(
	_ context.Context,
	request model.ListApplicationsRequest,
) ([]model.Application, int64, error) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	status := strings.ToLower(
		strings.TrimSpace(request.Status),
	)
	company := strings.ToLower(
		strings.TrimSpace(request.Company),
	)

	applications := make(
		[]model.Application,
		0,
		len(repository.applications),
	)

	for _, application := range repository.applications {
		if status != "" && application.Status != status {
			continue
		}

		if company != "" &&
			!strings.Contains(
				strings.ToLower(application.Company),
				company,
			) {
			continue
		}

		applications = append(
			applications,
			application,
		)
	}

	sort.Slice(
		applications,
		func(i, j int) bool {
			if applications[i].UpdatedAt.Equal(
				applications[j].UpdatedAt,
			) {
				return applications[i].ID >
					applications[j].ID
			}

			return applications[i].UpdatedAt.After(
				applications[j].UpdatedAt,
			)
		},
	)

	total := int64(len(applications))

	offset := (request.Page - 1) * request.PageSize

	if offset >= len(applications) {
		return []model.Application{}, total, nil
	}

	end := offset + request.PageSize
	if end > len(applications) {
		end = len(applications)
	}

	return applications[offset:end], total, nil

}

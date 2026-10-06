package repository

import (
	"context"
	"errors"
	"testing"

	"job-tracker-api/internal/model"
)

func TestMemoryApplicationRepositoryCreateAndGetByID(t *testing.T) {
	repository := NewMemoryApplicationRepository()

	ctx := context.Background()

	created, err := repository.Create(ctx, model.Application{
		Company:  "ByteDance",
		Position: "Backend Intern",
		Status:   "planned",
	})

	if err != nil {
		t.Fatalf("expected create to succeed, got error: %v", err)
	}
	if created.ID != 1 {
		t.Fatalf("expected ID 1, got %d", created.ID)
	}

	found, err := repository.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("expected application to exist, got error: %v", err)
	}

	if found != created {
		t.Fatalf("expected %+v, got %+v", created, found)
	}
}

func TestMemoryApplicationRepositoryDelete(t *testing.T) {
	repository := NewMemoryApplicationRepository()
	ctx := context.Background()

	created, err := repository.Create(ctx, model.Application{
		Company:  "Tencent",
		Position: "Backend Intern",
		Status:   "planned",
	})

	if err != nil {
		t.Fatalf("expected create to succeed, got error: %v", err)
	}

	if err := repository.Delete(ctx, created.ID); err != nil {
		t.Fatalf("expected delete to succeed, got error: %v", err)
	}

	_, err = repository.GetByID(ctx, created.ID)
	if !errors.Is(err, ErrApplicationNotFound) {
		t.Fatalf("expected ErrApplicationNotFound after delete, got: %v", err)
	}
}

func TestMemoryApplicationRepositoryDeleteNotFound(t *testing.T) {
	repository := NewMemoryApplicationRepository()
	ctx := context.Background()

	err := repository.Delete(ctx, 2)
	if !errors.Is(err, ErrApplicationNotFound) {
		t.Fatalf("expected ErrApplicationNotFound, got: %v", err)
	}
}

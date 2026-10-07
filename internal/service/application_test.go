package service

import (
	"context"
	"errors"
	"job-tracker-api/internal/model"
	"job-tracker-api/internal/repository"
	"strings"
	"testing"
)

type fakeApplicationRepository struct {
	created model.Application
}

func (repository *fakeApplicationRepository) Create(
	_ context.Context,
	application model.Application,
) (model.Application, error) {
	application.ID = 1
	repository.created = application
	return application, nil
}

func (repository *fakeApplicationRepository) GetByID(
	_ context.Context,
	id int64,
) (model.Application, error) {
	return model.Application{}, nil
}

func (repository *fakeApplicationRepository) Delete(
	_ context.Context,
	id int64,
) error {
	return nil
}

func TestApplicationServiceCreate(t *testing.T) {

	ctx := context.Background()

	repository := &fakeApplicationRepository{}
	applicationService := NewApplicationService(repository)

	appliedAt := "  2026-10-06  "

	created, err := applicationService.Create(ctx, model.CreateApplicationRequest{
		Company:   "  ByteDance  ",
		Position:  "  Backend Intern  ",
		Status:    "  INTERVIEW  ",
		AppliedAt: &appliedAt,
		Notes:     "  官网日常实习  ",
	})

	if err != nil {
		t.Fatalf("expected create to succeed, got error: %v", err)
	}

	if created.ID != 1 {
		t.Fatalf("expected ID 1, got %d", created.ID)
	}

	if created.Company != "ByteDance" {
		t.Fatalf("expected trimmed company, got %q", created.Company)
	}

	if created.Position != "Backend Intern" {
		t.Fatalf("expected trimmed position, got %q", created.Position)
	}

	if created.Status != "interview" {
		t.Fatalf("expected normalized status, got %q", created.Status)
	}
	if created.AppliedAt == nil {
		t.Fatal("expected applied_at, got nil")
	}

	if *created.AppliedAt != "2026-10-06" {
		t.Fatalf(
			"expected applied_at %q, got %q",
			"2026-10-06",
			*created.AppliedAt,
		)
	}
	if created.Notes != "官网日常实习" {
		t.Fatalf(
			"expected notes %q, got %q",
			"官网日常实习",
			created.Notes,
		)
	}
}

func TestApplicationServiceCreateValidation(t *testing.T) {
	ctx := context.Background()

	invalidAppliedAt := "2026-02-30"

	tests := []struct {
		name    string
		request model.CreateApplicationRequest
		wantErr error
	}{
		{
			name: "empty company",
			request: model.CreateApplicationRequest{
				Company:  "   ",
				Position: "Backend Intern",
				Status:   "planned",
			},
			wantErr: ErrCompanyRequired,
		},
		{
			name: "company too long",
			request: model.CreateApplicationRequest{
				Company:  strings.Repeat("中", 101),
				Position: "Backend Intern",
				Status:   "planned",
			},
			wantErr: ErrCompanyTooLong,
		},
		{
			name: "empty position",
			request: model.CreateApplicationRequest{
				Company:  "ByteDance",
				Position: "   ",
				Status:   "planned",
			},
			wantErr: ErrPositionRequired,
		},
		{
			name: "invalid status",
			request: model.CreateApplicationRequest{
				Company:  "ByteDance",
				Position: "Backend Intern",
				Status:   "unknown",
			},
			wantErr: ErrInvalidStatus,
		},
		{
			name: "invalid applied_at",
			request: model.CreateApplicationRequest{
				Company:   "ByteDance",
				Position:  "Backend Intern",
				Status:    "planned",
				AppliedAt: &invalidAppliedAt,
			},
			wantErr: ErrInvalidAppliedAt,
		},

		{
			name: "notes too long",
			request: model.CreateApplicationRequest{
				Company:  "ByteDance",
				Position: "Backend Intern",
				Status:   "planned",
				Notes:    strings.Repeat("中", 1001),
			},
			wantErr: ErrNotesTooLong,
		},
	}

	repository := &fakeApplicationRepository{}
	applicationService := NewApplicationService(repository)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := applicationService.Create(ctx, test.request)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf(
					"expected error %v, got %v",
					test.wantErr,
					err,
				)
			}
		})
	}
}

func TestApplicationServiceDeleteNotFound(t *testing.T) {

	ctx := context.Background()

	repo := repository.NewMemoryApplicationRepository()
	applicationService := NewApplicationService(repo)

	err := applicationService.Delete(ctx, 99)

	if !errors.Is(err, ErrApplicationNotFound) {
		t.Fatalf("expected ErrApplicationNotFound, got %v", err)
	}
}

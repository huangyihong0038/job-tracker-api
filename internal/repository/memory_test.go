package repository

import (
	"context"
	"errors"
	"testing"
	"time"

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

func TestMemoryApplicationRepositoryList(t *testing.T) {

	repository := NewMemoryApplicationRepository()
	ctx := context.Background()

	baseTime := time.Date(
		2026,
		time.October,
		7,
		10,
		0,
		0,
		0,
		time.UTC,
	)

	seedApplications := []model.Application{
		{
			Company:   "ByteDance",
			Position:  "Backend Intern",
			Status:    "planned",
			UpdatedAt: baseTime.Add(1 * time.Hour),
		},
		{
			Company:   "Tencent",
			Position:  "Backend Intern",
			Status:    "interview",
			UpdatedAt: baseTime.Add(3 * time.Hour),
		},
		{
			Company:   "ByteDance",
			Position:  "AI Application Intern",
			Status:    "interview",
			UpdatedAt: baseTime.Add(2 * time.Hour),
		},
		{
			Company:   "Meituan",
			Position:  "Backend Intern",
			Status:    "offer",
			UpdatedAt: baseTime.Add(3 * time.Hour),
		},
	}

	for _, application := range seedApplications {
		_, err := repository.Create(ctx, application)
		if err != nil {
			t.Fatalf(
				"seed application: %v",
				err,
			)
		}
	}

	tests := []struct {
		name      string
		request   model.ListApplicationsRequest
		wantIDs   []int64
		wantTotal int64
	}{
		{
			name: "first page",
			request: model.ListApplicationsRequest{
				Page:     1,
				PageSize: 2,
			},
			wantIDs:   []int64{4, 2},
			wantTotal: 4,
		},
		{
			name: "second page",
			request: model.ListApplicationsRequest{
				Page:     2,
				PageSize: 2,
			},
			wantIDs:   []int64{3, 1},
			wantTotal: 4,
		},
		{
			name: "filter by status",
			request: model.ListApplicationsRequest{
				Page:     1,
				PageSize: 20,
				Status:   "  INTERVIEW  ",
			},
			wantIDs:   []int64{2, 3},
			wantTotal: 2,
		},
		{
			name: "filter by company",
			request: model.ListApplicationsRequest{
				Page:     1,
				PageSize: 20,
				Company:  "  byte  ",
			},
			wantIDs:   []int64{3, 1},
			wantTotal: 2,
		},
		{
			name: "filter by status and company",
			request: model.ListApplicationsRequest{
				Page:     1,
				PageSize: 20,
				Status:   "interview",
				Company:  "ByteDance",
			},
			wantIDs:   []int64{3},
			wantTotal: 1,
		},
		{
			name: "page outside available data",
			request: model.ListApplicationsRequest{
				Page:     3,
				PageSize: 2,
			},
			wantIDs:   []int64{},
			wantTotal: 4,
		},
		{
			name: "no matching records",
			request: model.ListApplicationsRequest{
				Page:     1,
				PageSize: 20,
				Company:  "Nonexistent Company",
			},
			wantIDs:   []int64{},
			wantTotal: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			applications, total, err :=
				repository.List(ctx, test.request)

			if err != nil {
				t.Fatalf(
					"expected list to succeed, got: %v",
					err,
				)
			}

			if total != test.wantTotal {
				t.Fatalf(
					"expected total %d, got %d",
					test.wantTotal,
					total,
				)
			}

			if len(applications) != len(test.wantIDs) {
				t.Fatalf(
					"expected %d applications, got %d: %+v",
					len(test.wantIDs),
					len(applications),
					applications,
				)
			}

			if len(test.wantIDs) == 0 &&
				applications == nil {
				t.Fatal(
					"expected an empty non-nil slice",
				)
			}

			for index, wantID := range test.wantIDs {
				if applications[index].ID != wantID {
					t.Fatalf(
						"at index %d expected ID %d, got %d",
						index,
						wantID,
						applications[index].ID,
					)
				}
			}
		})
	}

}

func TestMemoryApplicationRepositoryUpdate(t *testing.T) {
	repository := NewMemoryApplicationRepository()
	ctx := context.Background()

	appliedAt := "2026-10-01"
	originalTime := time.Now().UTC().Add(
		-1 * time.Hour,
	)

	created, err := repository.Create(
		ctx,
		model.Application{
			Company:   "ByteDance",
			Position:  "Backend Intern",
			Status:    "planned",
			AppliedAt: &appliedAt,
			Notes:     "original notes",
			CreatedAt: originalTime,
			UpdatedAt: originalTime,
		},
	)

	if err != nil {
		t.Fatalf("create application: %v", err)
	}

	newStatus := "interview"
	newNotes := "updated notes"

	updated, err := repository.Update(
		ctx,
		created.ID,
		model.UpdateApplicationRequest{
			Status: model.OptionalString{
				Present: true,
				Value:   &newStatus,
			},
			Notes: model.OptionalString{
				Present: true,
				Value:   &newNotes,
			},
		},
	)

	if err != nil {
		t.Fatalf("update application: %v", err)
	}

	if updated.Status != newStatus {
		t.Fatalf(
			"expected status %q, got %q",
			newStatus,
			updated.Status,
		)
	}

	if updated.Notes != newNotes {
		t.Fatalf(
			"expected notes %q, got %q",
			newNotes,
			updated.Notes,
		)
	}

	if updated.Company != created.Company {
		t.Fatalf(
			"expected company to remain %q, got %q",
			created.Company,
			updated.Company,
		)
	}

	if updated.Position != created.Position {
		t.Fatalf(
			"expected position to remain %q, got %q",
			created.Position,
			updated.Position,
		)
	}

	if updated.AppliedAt == nil ||
		*updated.AppliedAt != appliedAt {
		t.Fatalf(
			"expected applied_at to remain %q, got %v",
			appliedAt,
			updated.AppliedAt,
		)
	}

	if !updated.CreatedAt.Equal(created.CreatedAt) {
		t.Fatalf(
			"expected created_at to remain unchanged",
		)
	}

	if !updated.UpdatedAt.After(created.UpdatedAt) {
		t.Fatalf(
			"expected updated_at after %v, got %v",
			created.UpdatedAt,
			updated.UpdatedAt,
		)
	}

	found, err := repository.GetByID(
		ctx,
		created.ID,
	)
	if err != nil {
		t.Fatalf("get updated application: %v", err)
	}

	if found.Status != newStatus ||
		found.Notes != newNotes {
		t.Fatalf(
			"updated application was not stored: %+v",
			found,
		)
	}
}

func TestMemoryApplicationRepositoryUpdateClearsAppliedAt(
	t *testing.T,
) {
	repository := NewMemoryApplicationRepository()
	ctx := context.Background()

	appliedAt := "2026-10-01"

	created, err := repository.Create(
		ctx,
		model.Application{
			Company:   "ByteDance",
			Position:  "Backend Intern",
			Status:    "applied",
			AppliedAt: &appliedAt,
		},
	)
	if err != nil {
		t.Fatalf("create application: %v", err)
	}

	updated, err := repository.Update(
		ctx,
		created.ID,
		model.UpdateApplicationRequest{
			AppliedAt: model.OptionalString{
				Present: true,
				Value:   nil,
			},
		},
	)
	if err != nil {
		t.Fatalf("update application: %v", err)
	}

	if updated.AppliedAt != nil {
		t.Fatalf(
			"expected applied_at to be nil, got %q",
			*updated.AppliedAt,
		)
	}
}

func TestMemoryApplicationRepositoryUpdateNotFound(
	t *testing.T,
) {
	repository := NewMemoryApplicationRepository()
	ctx := context.Background()

	status := "interview"

	_, err := repository.Update(
		ctx,
		999,
		model.UpdateApplicationRequest{
			Status: model.OptionalString{
				Present: true,
				Value:   &status,
			},
		},
	)

	if !errors.Is(err, ErrApplicationNotFound) {
		t.Fatalf(
			"expected ErrApplicationNotFound, got %v",
			err,
		)
	}
}

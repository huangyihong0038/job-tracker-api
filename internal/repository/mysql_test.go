package repository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"job-tracker-api/internal/config"
	"job-tracker-api/internal/model"
)

func TestMySQLApplicationRepositoryLifecycle(t *testing.T) {
	if os.Getenv("MYSQL_DSN") == "" {
		t.Skip("MYSQL_DSN is not set")
	}

	ctx := context.Background()
	db, err := config.OpenMySQL(ctx)
	if err != nil {
		t.Fatalf("connect to MySQL: %v", err)
	}
	defer db.Close()

	repo := NewMySQLApplicationRepository(db)
	appliedAt := "2026-10-07"

	created, err := repo.Create(ctx, model.Application{
		Company:   "Test Company",
		Position:  "Backend Intern",
		Status:    "planned",
		AppliedAt: &appliedAt,
		Notes:     "Created by MySQL repository integration test",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	if created.ID <= 0 {
		t.Fatalf("expected a positive ID, got %d", created.ID)
	}

	if created.AppliedAt == nil {
		t.Fatal("expected applied_at, got nil")
	}

	if *created.AppliedAt != appliedAt {
		t.Fatalf(
			"expected applied_at %q, got %q",
			appliedAt,
			*created.AppliedAt,
		)
	}

	if created.Notes != "Created by MySQL repository integration test" {
		t.Fatalf("unexpected notes: %q", created.Notes)
	}

	if created.CreatedAt.IsZero() {
		t.Fatal("expected created_at to be set")
	}

	if created.UpdatedAt.IsZero() {
		t.Fatal("expected updated_at to be set")
	}

	// 即使后面的断言失败，也清理这条测试数据。
	defer db.ExecContext(ctx, "DELETE FROM applications WHERE id = ?", created.ID)

	got, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("get by ID: %v", err)
	}
	if got.Company != created.Company {
		t.Fatalf(
			"expected company %q, got %q",
			created.Company,
			got.Company,
		)
	}

	if got.Position != created.Position {
		t.Fatalf(
			"expected position %q, got %q",
			created.Position,
			got.Position,
		)
	}

	if got.Status != created.Status {
		t.Fatalf(
			"expected status %q, got %q",
			created.Status,
			got.Status,
		)
	}

	if got.AppliedAt == nil {
		t.Fatal("expected retrieved applied_at, got nil")
	}

	if *got.AppliedAt != appliedAt {
		t.Fatalf(
			"expected retrieved applied_at %q, got %q",
			appliedAt,
			*got.AppliedAt,
		)
	}

	if got.Notes != created.Notes {
		t.Fatalf(
			"expected notes %q, got %q",
			created.Notes,
			got.Notes,
		)
	}

	if got.CreatedAt.IsZero() || got.UpdatedAt.IsZero() {
		t.Fatalf(
			"expected database timestamps, got created_at=%v updated_at=%v",
			got.CreatedAt,
			got.UpdatedAt,
		)
	}

	if err := repo.Delete(ctx, created.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}

	_, err = repo.GetByID(ctx, created.ID)
	if !errors.Is(err, ErrApplicationNotFound) {
		t.Fatalf("expected ErrApplicationNotFound after delete, got %v", err)
	}
}

func TestMySQLApplicationRepositoryNullAppliedAt(t *testing.T) {
	if os.Getenv("MYSQL_DSN") == "" {
		t.Skip("MYSQL_DSN is not set")
	}

	ctx := context.Background()

	db, err := config.OpenMySQL(ctx)
	if err != nil {
		t.Fatalf("connect to MySQL: %v", err)
	}
	defer db.Close()

	repo := NewMySQLApplicationRepository(db)

	created, err := repo.Create(ctx, model.Application{
		Company:  "Null Date Test Company",
		Position: "Backend Intern",
		Status:   "planned",
		Notes:    "",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	// 即使后续断言失败，也删除测试数据。
	defer db.ExecContext(
		ctx,
		"DELETE FROM applications WHERE id = ?",
		created.ID,
	)

	if created.ID <= 0 {
		t.Fatalf(
			"expected a positive ID, got %d",
			created.ID,
		)
	}

	if created.AppliedAt != nil {
		t.Fatalf(
			"expected nil applied_at, got %q",
			*created.AppliedAt,
		)
	}

	if created.Notes != "" {
		t.Fatalf(
			"expected empty notes, got %q",
			created.Notes,
		)
	}

	got, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("get by ID: %v", err)
	}

	if got.AppliedAt != nil {
		t.Fatalf(
			"expected retrieved applied_at to be nil, got %q",
			*got.AppliedAt,
		)
	}

	if got.Notes != "" {
		t.Fatalf(
			"expected retrieved notes to be empty, got %q",
			got.Notes,
		)
	}

	if got.CreatedAt.IsZero() {
		t.Fatal("expected created_at to be set")
	}

	if got.UpdatedAt.IsZero() {
		t.Fatal("expected updated_at to be set")
	}
}

func TestMySQLApplicationRepositoryList(t *testing.T) {
	if os.Getenv("MYSQL_DSN") == "" {
		t.Skip("MYSQL_DSN is not set")
	}

	ctx := context.Background()

	db, err := config.OpenMySQL(ctx)
	if err != nil {
		t.Fatalf("connect to MySQL: %v", err)
	}
	defer db.Close()

	repo := NewMySQLApplicationRepository(db)

	marker := fmt.Sprintf(
		"ListTest-%d",
		time.Now().UnixNano(),
	)

	seedApplications := []model.Application{
		{
			Company:  marker + "-ByteDance",
			Position: "Backend Intern",
			Status:   "planned",
		},
		{
			Company:  marker + "-Tencent",
			Position: "Backend Intern",
			Status:   "interview",
		},
		{
			Company:  marker + "-ByteDance-AI",
			Position: "AI Application Intern",
			Status:   "interview",
		},
		{
			Company:  marker + "-Meituan",
			Position: "Backend Intern",
			Status:   "offer",
		},
	}

	created := make(
		[]model.Application,
		0,
		len(seedApplications),
	)

	defer func() {
		for _, application := range created {
			_, cleanupErr := db.ExecContext(
				ctx,
				"DELETE FROM applications WHERE id = ?",
				application.ID,
			)
			if cleanupErr != nil {
				t.Errorf(
					"cleanup application %d: %v",
					application.ID,
					cleanupErr,
				)
			}
		}
	}()

	for _, application := range seedApplications {
		result, createErr :=
			repo.Create(ctx, application)
		if createErr != nil {
			t.Fatalf(
				"seed application: %v",
				createErr,
			)
		}

		created = append(created, result)
	}

	allApplications, allTotal, err := repo.List(
		ctx,
		model.ListApplicationsRequest{
			Page:     1,
			PageSize: 1,
		},
	)
	if err != nil {
		t.Fatalf(
			"list without filters: %v",
			err,
		)
	}

	if len(allApplications) != 1 {
		t.Fatalf(
			"expected one application on first page, got %d",
			len(allApplications),
		)
	}

	if allTotal < int64(len(created)) {
		t.Fatalf(
			"expected total at least %d, got %d",
			len(created),
			allTotal,
		)
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
				Company:  marker,
			},
			wantIDs: []int64{
				created[3].ID,
				created[2].ID,
			},
			wantTotal: 4,
		},
		{
			name: "second page",
			request: model.ListApplicationsRequest{
				Page:     2,
				PageSize: 2,
				Company:  marker,
			},
			wantIDs: []int64{
				created[1].ID,
				created[0].ID,
			},
			wantTotal: 4,
		},
		{
			name: "filter by status",
			request: model.ListApplicationsRequest{
				Page:     1,
				PageSize: 20,
				Status:   "interview",
				Company:  marker,
			},
			wantIDs: []int64{
				created[2].ID,
				created[1].ID,
			},
			wantTotal: 2,
		},
		{
			name: "filter by company",
			request: model.ListApplicationsRequest{
				Page:     1,
				PageSize: 20,
				Company:  marker + "-Byte",
			},
			wantIDs: []int64{
				created[2].ID,
				created[0].ID,
			},
			wantTotal: 2,
		},
		{
			name: "no matching records",
			request: model.ListApplicationsRequest{
				Page:     1,
				PageSize: 20,
				Company:  marker + "-Missing",
			},
			wantIDs:   []int64{},
			wantTotal: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			applications, total, listErr :=
				repo.List(ctx, test.request)

			if listErr != nil {
				t.Fatalf(
					"list applications: %v",
					listErr,
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

func TestMySQLApplicationRepositoryUpdate(t *testing.T) {

	if os.Getenv("MYSQL_DSN") == "" {
		t.Skip("MYSQL_DSN is not set")
	}

	ctx := context.Background()

	db, err := config.OpenMySQL(ctx)
	if err != nil {
		t.Fatalf("connect to MySQL: %v", err)
	}
	defer db.Close()

	repo := NewMySQLApplicationRepository(db)

	appliedAt := "2026-10-08"

	created, err := repo.Create(
		ctx,
		model.Application{
			Company:   "Update Test Company",
			Position:  "Backend Intern",
			Status:    "planned",
			AppliedAt: &appliedAt,
			Notes:     "original notes",
		},
	)

	if err != nil {
		t.Fatalf("create application: %v", err)
	}

	defer func() {
		_, cleanupErr := db.ExecContext(
			ctx,
			"DELETE FROM applications WHERE id = ?",
			created.ID,
		)
		if cleanupErr != nil {
			t.Errorf(
				"cleanup application %d: %v",
				created.ID,
				cleanupErr,
			)
		}
	}()

	newStatus := "interview"
	newNotes := "updated notes"

	updated, err := repo.Update(
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

	cleared, err := repo.Update(
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
		t.Fatalf("clear applied_at: %v", err)
	}

	if cleared.AppliedAt != nil {
		t.Fatalf(
			"expected applied_at to be nil, got %q",
			*cleared.AppliedAt,
		)
	}

	if err := repo.Delete(ctx, created.ID); err != nil {
		t.Fatalf("delete application: %v", err)
	}

	_, err = repo.Update(
		ctx,
		created.ID,
		model.UpdateApplicationRequest{
			Status: model.OptionalString{
				Present: true,
				Value:   &newStatus,
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

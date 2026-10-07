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

	listFunc func(
		model.ListApplicationsRequest,
	) ([]model.Application, int64, error)
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

func (repository *fakeApplicationRepository) List(
	_ context.Context,
	request model.ListApplicationsRequest,
) ([]model.Application, int64, error) {
	if repository.listFunc != nil {
		return repository.listFunc(request)
	}

	return []model.Application{}, 0, nil
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

func TestApplicationServiceList(t *testing.T) {
	ctx := context.Background()

	repository := &fakeApplicationRepository{
		listFunc: func(
			request model.ListApplicationsRequest,
		) ([]model.Application, int64, error) {
			if request.Page != 2 {
				t.Fatalf(
					"expected page 2, got %d",
					request.Page,
				)
			}

			if request.PageSize != 10 {
				t.Fatalf(
					"expected page size 10, got %d",
					request.PageSize,
				)
			}

			if request.Status != "interview" {
				t.Fatalf(
					"expected normalized status interview, got %q",
					request.Status,
				)
			}

			if request.Company != "Byte" {
				t.Fatalf(
					"expected trimmed company Byte, got %q",
					request.Company,
				)
			}

			return []model.Application{
				{
					ID:      2,
					Company: "ByteDance",
					Status:  "interview",
				},
				{
					ID:      1,
					Company: "ByteDance",
					Status:  "interview",
				},
			}, 25, nil
		},
	}

	applicationService :=
		NewApplicationService(repository)

	response, err := applicationService.List(
		ctx,
		model.ListApplicationsRequest{
			Page:     2,
			PageSize: 10,
			Status:   "  INTERVIEW  ",
			Company:  "  Byte  ",
		},
	)
	if err != nil {
		t.Fatalf(
			"expected list to succeed, got: %v",
			err,
		)
	}

	if len(response.Data) != 2 {
		t.Fatalf(
			"expected 2 applications, got %d",
			len(response.Data),
		)
	}

	if response.Data[0].ID != 2 ||
		response.Data[1].ID != 1 {
		t.Fatalf(
			"unexpected application order: %+v",
			response.Data,
		)
	}

	if response.Pagination.Page != 2 {
		t.Fatalf(
			"expected page 2, got %d",
			response.Pagination.Page,
		)
	}

	if response.Pagination.PageSize != 10 {
		t.Fatalf(
			"expected page size 10, got %d",
			response.Pagination.PageSize,
		)
	}

	if response.Pagination.Total != 25 {
		t.Fatalf(
			"expected total 25, got %d",
			response.Pagination.Total,
		)
	}
}

func TestApplicationServiceListValidation(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		request model.ListApplicationsRequest
		wantErr error
	}{
		{
			name: "page is zero",
			request: model.ListApplicationsRequest{
				Page:     0,
				PageSize: 20,
			},
			wantErr: ErrInvalidPage,
		},
		{
			name: "page size is zero",
			request: model.ListApplicationsRequest{
				Page:     1,
				PageSize: 0,
			},
			wantErr: ErrInvalidPageSize,
		},
		{
			name: "page size exceeds maximum",
			request: model.ListApplicationsRequest{
				Page:     1,
				PageSize: 101,
			},
			wantErr: ErrInvalidPageSize,
		},
		{
			name: "invalid status",
			request: model.ListApplicationsRequest{
				Page:     1,
				PageSize: 20,
				Status:   "unknown",
			},
			wantErr: ErrInvalidStatus,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repositoryCalled := false

			repository := &fakeApplicationRepository{
				listFunc: func(
					request model.ListApplicationsRequest,
				) ([]model.Application, int64, error) {
					repositoryCalled = true
					return nil, 0, nil
				},
			}

			applicationService :=
				NewApplicationService(repository)

			_, err := applicationService.List(
				ctx,
				test.request,
			)

			if !errors.Is(err, test.wantErr) {
				t.Fatalf(
					"expected error %v, got %v",
					test.wantErr,
					err,
				)
			}

			if repositoryCalled {
				t.Fatal(
					"repository should not be called for invalid request",
				)
			}
		})
	}
}

func TestApplicationServiceListRepositoryError(t *testing.T) {
	ctx := context.Background()

	repositoryError :=
		errors.New("database connection failed")

	repository := &fakeApplicationRepository{
		listFunc: func(
			request model.ListApplicationsRequest,
		) ([]model.Application, int64, error) {
			return nil, 0, repositoryError
		},
	}

	applicationService :=
		NewApplicationService(repository)

	_, err := applicationService.List(
		ctx,
		model.ListApplicationsRequest{
			Page:     1,
			PageSize: 20,
		},
	)

	if !errors.Is(err, repositoryError) {
		t.Fatalf(
			"expected repository error %v, got %v",
			repositoryError,
			err,
		)
	}
}

func TestApplicationServiceListNormalizesNilData(t *testing.T) {
	ctx := context.Background()

	repository := &fakeApplicationRepository{
		listFunc: func(
			request model.ListApplicationsRequest,
		) ([]model.Application, int64, error) {
			return nil, 0, nil
		},
	}

	applicationService :=
		NewApplicationService(repository)

	response, err := applicationService.List(
		ctx,
		model.ListApplicationsRequest{
			Page:     1,
			PageSize: 20,
		},
	)
	if err != nil {
		t.Fatalf(
			"expected list to succeed, got: %v",
			err,
		)
	}

	if response.Data == nil {
		t.Fatal(
			"expected empty non-nil data slice",
		)
	}

	if len(response.Data) != 0 {
		t.Fatalf(
			"expected empty data, got %+v",
			response.Data,
		)
	}
}

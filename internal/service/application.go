package service

import (
	"context"
	"errors"
	"job-tracker-api/internal/model"
	"job-tracker-api/internal/repository"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	ErrCompanyTooLong      = errors.New("company is too long")
	ErrCompanyRequired     = errors.New("company is required")
	ErrPositionRequired    = errors.New("position is required")
	ErrPositionTooLong     = errors.New("position is too long")
	ErrInvalidStatus       = errors.New("invalid status")
	ErrApplicationNotFound = errors.New("application not found")
	ErrInvalidAppliedAt    = errors.New("applied_at must be in YYYY-MM-DD format")
	ErrNotesTooLong        = errors.New("notes is too long")
)

type ApplicationRepository interface {
	Create(context.Context, model.Application) (model.Application, error)
	GetByID(context.Context, int64) (model.Application, error)
	Delete(context.Context, int64) error
}

type ApplicationService struct {
	repository ApplicationRepository
}

func NewApplicationService(repo ApplicationRepository) *ApplicationService {
	return &ApplicationService{
		repository: repo,
	}

}

func ValidStatus(
	status string,
) error {
	switch status {
	case "planned",
		"applied",
		"written_test",
		"interview",
		"offer",
		"rejected",
		"withdrawn":
		return nil
	default:
		return ErrInvalidStatus
	}
}

func (service *ApplicationService) Create(
	ctx context.Context,
	request model.CreateApplicationRequest,
) (model.Application, error) {
	company := strings.TrimSpace(request.Company)
	position := strings.TrimSpace(request.Position)
	status := strings.ToLower(strings.TrimSpace(request.Status))
	notes := strings.TrimSpace(request.Notes)

	var appliedAt *string
	if request.AppliedAt != nil {

		value := strings.TrimSpace(*request.AppliedAt)

		if _, err := time.Parse("2006-01-02", value); err != nil {

			return model.Application{}, ErrInvalidAppliedAt
		}

		appliedAt = &value
	}

	if company == "" {
		return model.Application{}, ErrCompanyRequired
	}

	if utf8.RuneCountInString(company) > 100 {
		return model.Application{}, ErrCompanyTooLong
	}

	if position == "" {
		return model.Application{}, ErrPositionRequired
	}

	if utf8.RuneCountInString(position) > 100 {
		return model.Application{}, ErrPositionTooLong
	}

	if err := ValidStatus(status); err != nil {
		return model.Application{}, ErrInvalidStatus
	}

	if utf8.RuneCountInString(notes) > 1000 {
		return model.Application{}, ErrNotesTooLong
	}

	application := model.Application{
		Company:   company,
		Position:  position,
		Status:    status,
		AppliedAt: appliedAt,
		Notes:     notes,
	}

	return service.repository.Create(ctx, application)
}

func (service *ApplicationService) GetByID(
	ctx context.Context,
	id int64,
) (model.Application, error) {

	application, err := service.repository.GetByID(ctx, id)

	if err != nil {
		if errors.Is(err, repository.ErrApplicationNotFound) {
			return model.Application{}, ErrApplicationNotFound
		}

		return model.Application{}, err
	}

	return application, nil
}

func (service *ApplicationService) Delete(
	ctx context.Context,
	id int64,
) error {
	err := service.repository.Delete(ctx, id)
	if errors.Is(err, repository.ErrApplicationNotFound) {
		return ErrApplicationNotFound
	}
	return err
}

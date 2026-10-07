package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"job-tracker-api/internal/model"
)

type MySQLApplicationRepository struct {
	db *sql.DB
}

func NewMySQLApplicationRepository(
	db *sql.DB,
) *MySQLApplicationRepository {

	return &MySQLApplicationRepository{db: db}
}

func (r *MySQLApplicationRepository) Create(
	ctx context.Context,
	application model.Application,
) (model.Application, error) {

	var appliedAt any

	if application.AppliedAt != nil {
		appliedAt = *application.AppliedAt
	}

	result, err := r.db.ExecContext(
		ctx,
		`INSERT INTO applications (company, position, status,applied_at,notes)
		 VALUES (?, ?, ?, ?, ?)`,
		application.Company,
		application.Position,
		application.Status,
		appliedAt,
		application.Notes,
	)

	if err != nil {
		return model.Application{}, fmt.Errorf("insert application: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return model.Application{}, fmt.Errorf("get inserted ID: %w", err)
	}

	created, err := r.GetByID(ctx, id)
	if err != nil {
		return model.Application{},
			fmt.Errorf("get created application: %w", err)
	}

	return created, nil
}

func (r *MySQLApplicationRepository) GetByID(
	ctx context.Context,
	id int64,
) (model.Application, error) {
	var application model.Application
	var appliedAt sql.NullTime

	err := r.db.QueryRowContext(
		ctx,
		`SELECT
			id,
			company,
			position,
			status,
			applied_at,
			notes,
			created_at,
			updated_at
		 FROM applications
		 WHERE id = ?`,
		id,
	).Scan(
		&application.ID,
		&application.Company,
		&application.Position,
		&application.Status,
		&appliedAt,
		&application.Notes,
		&application.CreatedAt,
		&application.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return model.Application{}, ErrApplicationNotFound
	}
	if err != nil {
		return model.Application{}, fmt.Errorf("get application by ID: %w", err)
	}

	if appliedAt.Valid {
		value := appliedAt.Time.Format("2006-01-02")
		application.AppliedAt = &value
	}

	return application, nil
}

func (r *MySQLApplicationRepository) Delete(
	ctx context.Context,
	id int64,
) error {
	result, err := r.db.ExecContext(
		ctx,
		`DELETE FROM applications WHERE id = ?`,
		id,
	)
	if err != nil {
		return fmt.Errorf("delete application: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get deleted row count: %w", err)
	}
	if affected == 0 {
		return ErrApplicationNotFound
	}

	return nil
}

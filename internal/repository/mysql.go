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
	result, err := r.db.ExecContext(
		ctx,
		`INSERT INTO applications (company, position, status)
		 VALUES (?, ?, ?)`,
		application.Company,
		application.Position,
		application.Status,
	)

	if err != nil {
		return model.Application{}, fmt.Errorf("insert application: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return model.Application{}, fmt.Errorf("get inserted ID: %w", err)
	}

	application.ID = id
	return application, nil
}

func (r *MySQLApplicationRepository) GetByID(
	ctx context.Context,
	id int64,
) (model.Application, error) {
	var application model.Application

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, company, position, status
		 FROM applications
		 WHERE id = ?`,
		id,
	).Scan(
		&application.ID,
		&application.Company,
		&application.Position,
		&application.Status,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return model.Application{}, ErrApplicationNotFound
	}
	if err != nil {
		return model.Application{}, fmt.Errorf("get application by ID: %w", err)
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

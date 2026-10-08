package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

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

func (r *MySQLApplicationRepository) List(
	ctx context.Context,
	request model.ListApplicationsRequest,
) ([]model.Application, int64, error) {

	conditions := make([]string, 0, 2)
	args := make([]any, 0, 2)

	if request.Status != "" {
		conditions = append(
			conditions,
			"status = ?",
		)

		args = append(
			args,
			request.Status,
		)
	}

	if request.Company != "" {
		conditions = append(
			conditions,
			"company LIKE ?",
		)
		args = append(
			args,
			"%"+request.Company+"%",
		)
	}

	whereClause := ""

	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	var total int64

	countQuery :=
		"SELECT COUNT(*) FROM applications" + whereClause

	err := r.db.QueryRowContext(
		ctx,
		countQuery,
		args...,
	).Scan(&total)
	if err != nil {
		return nil, 0,
			fmt.Errorf("count applications: %w", err)
	}

	query := `SELECT
		id,
		company,
		position,
		status,
		applied_at,
		notes,
		created_at,
		updated_at
	FROM applications` +
		whereClause +
		` ORDER BY updated_at DESC, id DESC
		LIMIT ? OFFSET ?`

	offset := (request.Page - 1) * request.PageSize

	queryArgs := append([]any{}, args...)

	queryArgs = append(
		queryArgs,
		request.PageSize,
		offset,
	)

	rows, err := r.db.QueryContext(
		ctx,
		query,
		queryArgs...,
	)
	if err != nil {
		return nil, 0,
			fmt.Errorf("list applications: %w", err)
	}
	defer rows.Close()

	applications := make(
		[]model.Application,
		0,
	)

	for rows.Next() {
		var application model.Application
		var appliedAt sql.NullTime

		err := rows.Scan(
			&application.ID,
			&application.Company,
			&application.Position,
			&application.Status,
			&appliedAt,
			&application.Notes,
			&application.CreatedAt,
			&application.UpdatedAt,
		)
		if err != nil {
			return nil, 0,
				fmt.Errorf("scan application: %w", err)
		}

		if appliedAt.Valid {
			value :=
				appliedAt.Time.Format("2006-01-02")
			application.AppliedAt = &value
		}

		applications = append(
			applications,
			application,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, 0,
			fmt.Errorf("iterate applications: %w", err)
	}

	return applications, total, nil

}

func (r *MySQLApplicationRepository) Update(
	ctx context.Context,
	id int64,
	request model.UpdateApplicationRequest,
) (model.Application, error) {
	assignments := make([]string, 0, 5)
	args := make([]any, 0, 6)

	if request.Company.Present &&
		request.Company.Value != nil {
		assignments = append(
			assignments,
			"company = ?",
		)
		args = append(
			args,
			*request.Company.Value,
		)
	}

	if request.Position.Present &&
		request.Position.Value != nil {
		assignments = append(
			assignments,
			"position = ?",
		)
		args = append(
			args,
			*request.Position.Value,
		)
	}

	if request.Status.Present &&
		request.Status.Value != nil {
		assignments = append(
			assignments,
			"status = ?",
		)
		args = append(
			args,
			*request.Status.Value,
		)
	}

	if request.AppliedAt.Present {
		assignments = append(
			assignments,
			"applied_at = ?",
		)

		if request.AppliedAt.Value == nil {
			args = append(args, nil)
		} else {
			args = append(
				args,
				*request.AppliedAt.Value,
			)
		}
	}

	if request.Notes.Present &&
		request.Notes.Value != nil {
		assignments = append(
			assignments,
			"notes = ?",
		)
		args = append(
			args,
			*request.Notes.Value,
		)
	}

	if len(assignments) == 0 {
		return r.GetByID(ctx, id)
	}

	args = append(args, id)

	query := "UPDATE applications SET " +
		strings.Join(assignments, ", ") +
		" WHERE id = ?"

	_, err := r.db.ExecContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return model.Application{},
			fmt.Errorf(
				"update application: %w",
				err,
			)
	}

	updated, err := r.GetByID(ctx, id)
	if err != nil {
		if errors.Is(
			err,
			ErrApplicationNotFound,
		) {
			return model.Application{},
				ErrApplicationNotFound
		}

		return model.Application{},
			fmt.Errorf(
				"get updated application: %w",
				err,
			)
	}

	return updated, nil
}

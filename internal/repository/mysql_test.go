package repository

import (
	"context"
	"errors"
	"os"
	"testing"

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

	created, err := repo.Create(ctx, model.Application{
		Company:  "Test Company",
		Position: "Backend Intern",
		Status:   "planned",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	// 即使后面的断言失败，也清理这条测试数据。
	defer db.ExecContext(ctx, "DELETE FROM applications WHERE id = ?", created.ID)

	got, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("get by ID: %v", err)
	}
	if got.Company != created.Company || got.Status != created.Status {
		t.Fatalf("unexpected application: %+v", got)
	}

	if err := repo.Delete(ctx, created.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}

	_, err = repo.GetByID(ctx, created.ID)
	if !errors.Is(err, ErrApplicationNotFound) {
		t.Fatalf("expected ErrApplicationNotFound after delete, got %v", err)
	}
}

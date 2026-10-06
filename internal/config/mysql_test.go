package config

import (
	"context"
	"os"
	"testing"
)

func TestOpenMySQLIntegration(t *testing.T) {
	if os.Getenv("MYSQL_DSN") == "" {
		t.Skip("MYSQL_DSN is not set")
	}

	db, err := OpenMySQL(context.Background())

	if err != nil {

		t.Fatalf("connect to MySQL: %v", err)
	}
	defer db.Close()

	var databaseName string

	if err := db.QueryRowContext(
		context.Background(),
		"SELECT DATABASE()",
	).Scan(&databaseName); err != nil {
		t.Fatalf("query database name: %v", err)
	}

	if databaseName != "job_tracker" {
		t.Fatalf("expected job_tracker, got %q", databaseName)
	}

}

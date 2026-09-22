package migrations

import (
	"database/sql"
	"fmt"
	"os"
)

const migrationName = "001_create_credit_cards_and_transactions"

func Run(db *sql.DB, path string) error {
	script, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read migration file %q: %w", path, err)
	}

	if _, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`); err != nil {
		return fmt.Errorf("create schema_migrations table: %w", err)
	}

	var applied bool
	if err = db.QueryRow(
		"SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE version = $1)",
		migrationName,
	).Scan(&applied); err != nil {
		return fmt.Errorf("check migration %q: %w", migrationName, err)
	}
	if applied {
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin migration %q: %w", migrationName, err)
	}
	defer tx.Rollback()

	if _, err = tx.Exec(string(script)); err != nil {
		return fmt.Errorf("execute migration %q: %w", migrationName, err)
	}
	if _, err = tx.Exec(
		"INSERT INTO schema_migrations (version) VALUES ($1)",
		migrationName,
	); err != nil {
		return fmt.Errorf("record migration %q: %w", migrationName, err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit migration %q: %w", migrationName, err)
	}
	return nil
}

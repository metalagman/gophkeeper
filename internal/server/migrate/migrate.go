package migrate

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func init() {
	goose.SetBaseFS(embedMigrations)
	_ = goose.SetDialect("postgres")
}

func Up(db *sql.DB) error {
	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}
	return nil
}

func Down(db *sql.DB) error {
	if err := goose.Down(db, "migrations"); err != nil {
		return fmt.Errorf("goose down: %w", err)
	}
	return nil
}

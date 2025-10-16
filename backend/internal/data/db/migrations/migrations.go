// Package migrations handles the database migrations using goose and embedded sql files.
package migrations

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/stdlib"

	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
)

// AdvisoryLock is the advisory lock ID used to prevent multiple migrations
// from running concurrently. This integer was randomly selected and has no
// inherent meaning.
var AdvisoryLock = 103945

//go:embed sql
var migrationFS embed.FS

// The goose library uses globals to store configuration state, so we run these once in the init.
// Similarly to [loggerOnce] this exists to allow concurrent test execution, and has no real impact
// in production.
func init() { // nolint:gochecknoinits
	err := goose.SetDialect("pgx")
	if err != nil {
		panic(fmt.Errorf("failed to set dialect: %w", err))
	}
	goose.SetBaseFS(migrationFS)
}

func SetLogger(log zerolog.Logger) {
	goose.SetLogger(logger{log: log})
}

func Up(log zerolog.Logger, db *sql.DB) error {
	if err := goose.Up(db, "sql"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Info().Msg("successfully ran migrations")
	return nil
}

func DownByOne(ctx context.Context, log zerolog.Logger, db *sql.DB) error {
	if err := goose.DownContext(ctx, db, "sql"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Info().Msg("successfully ran migrations")
	return nil
}

func Rollback(log zerolog.Logger, db *sql.DB) error {
	if err := goose.Down(db, "sql"); err != nil {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	log.Info().Msg("successfully rolled back migrations")
	return nil
}

var _ goose.Logger = &logger{}

type logger struct {
	log zerolog.Logger
}

// Fatalf implements goose.Logger.
func (l logger) Fatalf(format string, v ...any) {
	trimmed := strings.TrimSuffix(format, "\n")
	l.log.Fatal().Msgf(trimmed, v...)
}

// Printf implements goose.Logger.
func (l logger) Printf(format string, v ...any) {
	trimmed := strings.TrimSuffix(format, "\n")
	l.log.Info().Msgf(trimmed, v...)
}

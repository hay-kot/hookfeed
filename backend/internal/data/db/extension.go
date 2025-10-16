package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/hay-kot/hookfeed/backend/internal/data/db/migrations"
	"github.com/hay-kot/httpkit/errtrace"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

// QueriesExt is an extension of the generated Queries struct which also depends directly on the internal
// sql connection and allows for easier transaction handling and some basic utility methods for working
// with the database.
type QueriesExt struct {
	*Queries
	conn *pgxpool.Pool
}

// Close closes the connection.
func (qe *QueriesExt) Close(ctx context.Context) error {
	qe.conn.Close()
	return nil
}

// Ctx can be used to check the context for a keyed transaction. This can be used
// to cordinate transactions across stores.
//
// If no transaction exists in the context, the existing object is returned.
func (qe *QueriesExt) Ctx(ctx context.Context) *QueriesExt {
	tx, ok := getTransaction(ctx)
	if !ok {
		return qe
	}

	return &QueriesExt{qe.WithTx(tx), qe.conn}
}

// WithinTx runs the given function in a transaction.
func (qe *QueriesExt) WithinTx(ctx context.Context, fn func(*QueriesExt) error) error {
	ctx, tx, err := WithTransaction(ctx, qe)
	if err != nil {
		return err
	}

	qext := &QueriesExt{qe.WithTx(tx), qe.conn}
	if err := fn(qext); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func NewExt(ctx context.Context, logger zerolog.Logger, config Config, runMigrations bool) (*QueriesExt, error) {
	var conn *pgxpool.Pool
	var err error

	var (
		retries = 5
		wait    = 1
		dsn     = config.DSN()
	)

	for {
		conn, err = pgxpool.New(ctx, dsn)
		if err == nil {
			err = conn.Ping(ctx)
			if err == nil {
				break
			}
		}

		if retries == 0 {
			return nil, err
		}

		retries--
		logger.Warn().Err(err).Int("retries", retries).Msg("failed to ping database, retrying...")
		time.Sleep(time.Duration(wait) * time.Second)
		wait *= 2
	}

	if runMigrations {
		var (
			retries = 5
			wait    = 5
		)

		for {
			_, err = conn.Exec(ctx, "SELECT pg_advisory_lock($1)", migrations.AdvisoryLock)
			if err != nil {
				if retries == 0 {
					return nil, err
				}

				retries--
				logger.Warn().Err(err).Int("retries", retries).Msg("failed to obtain advisory lock, retrying...")
				time.Sleep(time.Duration(wait) * time.Second)
				wait *= 2
				continue
			}

			logger.Info().Msg("obtained advisory lock for migrations")
			defer func() {
				_, err = conn.Exec(ctx, "SELECT pg_advisory_unlock($1)", migrations.AdvisoryLock)
				if err != nil {
					logger.Error().Err(err).Msg("failed to release advisory")
				}
			}()

			break
		}

		stdlibConn, err := sql.Open("pgx", dsn)
		if err != nil {
			return nil, err
		}
		defer func() {
			if err := stdlibConn.Close(); err != nil {
				logger.Error().Err(err).Msg("failed to close stdlib connection")
			}
		}()

		err = migrations.Up(logger, stdlibConn)
		if err != nil {
			return nil, err
		}
	}

	return &QueriesExt{
		Queries: New(conn),
		conn:    conn,
	}, nil
}

const (
	OrderByAsc  = "asc"
	OrderByDesc = "desc"
)

// OrderBy is a utility function to 'compile' an order by clause used for dynamic ordering in sqlc
// in practive this means simple concatenation `column + ':' + direction`. We use this method to
// keep consistency and validate direction with fallback to defaults.
//
// This method does not validate the column to order by. See [Queries.AdminGetUsersWithCounts]
// usage for an example query for how this can be used.
func OrderBy(column, direction string) (string, error) {
	if column == "" {
		return "", errtrace.New("empty column name for order by clause")
	}

	dir := ""
	switch direction {
	case "asc", "ASC", "ascending":
		dir = "asc"
	case "desc", "DESC", "descending":
		dir = "desc"
	default:
		return "", errtrace.New("invalid direction '%s' for order by clause", direction)
	}

	return column + ":" + dir, nil
}

// OrderByWithDefaults extends OrderBy to provide a default ordering if the column is empty.
// or if the direction is invalid. This is useful for providing a default ordering and column
// when working with optional/pointer types.
func OrderByWithDefaults(inputColumn *string, defaultColumn string, inputDirection *string, defaultDirection string) (string, error) {
	column := defaultColumn
	if inputColumn != nil && *inputColumn != "" {
		column = *inputColumn
	}

	dir := defaultDirection
	if inputDirection != nil && *inputDirection != "" {
		dir = *inputDirection
	}

	return OrderBy(column, dir)
}

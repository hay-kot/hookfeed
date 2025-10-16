package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type queriesEmbeddedTx string

const queriesEmbeddedTxKey = queriesEmbeddedTx("queries:tx")

// WithTransaction creates a new transaction using the underlying database connection
// and attached it to the context, which can be used to compose transactions across
// domains using the same idomatic APIs.
//
// When composing a transactions like this you must call tx.Rollback or tx.Commit.
func WithTransaction(ctx context.Context, ext *QueriesExt) (context.Context, pgx.Tx, error) {
	tx, err := ext.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, nil, err
	}

	ctx = context.WithValue(ctx, queriesEmbeddedTxKey, tx)
	return ctx, tx, nil
}

func getTransaction(ctx context.Context) (tx pgx.Tx, ok bool) {
	tx, ok = ctx.Value(queriesEmbeddedTxKey).(pgx.Tx)
	return tx, ok
}

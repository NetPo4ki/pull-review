package repo

import (
	"context"

	"github.com/NetPo4ki/pull-review/internal/repo/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TxManager struct{ pool *pgxpool.Pool }

func NewTxManager(pool *pgxpool.Pool) *TxManager { return &TxManager{pool: pool} }

func (tm *TxManager) Do(ctx context.Context, fn func(ctx context.Context, q *sqlc.Queries) error) error {
	tx, err := tm.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	q := sqlc.New(tm.pool).WithTx(tx)
	if err := fn(ctx, q); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

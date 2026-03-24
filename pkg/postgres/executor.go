package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type TxKey struct{}

type Executor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

func (p *Postgres) GetExecutor(ctx context.Context) Executor {
	if tx, ok := ctx.Value(TxKey{}).(pgx.Tx); ok {
		return tx
	}

	return p.Pool
}

func (p *Postgres) WithinTransaction(ctx context.Context, f func(ctx context.Context) error) error {
	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Postgres - WithinTransaction - p.Pool.Begin: %w", err)
	}

	err = f(context.WithValue(ctx, TxKey{}, tx))
	if err != nil {
		_ = tx.Rollback(ctx)

		return fmt.Errorf("Postgres - WithinTransaction: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("Postgres - WithinTransaction - tx.Commit: %w", err)
	}

	return nil
}

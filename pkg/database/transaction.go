package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type TransactionManager interface {
	BeginTransaction(ctx context.Context) (context.Context, error)
	CommitTransaction(ctx context.Context) error
	RollbackTransaction(ctx context.Context) error
}

func (pg *PgContext) BeginTransaction(ctx context.Context) (context.Context, error) {
	tx, err := pg.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	ctx = context.WithValue(ctx, transactionKey, tx)

	return ctx, nil
}

func (pg *PgContext) CommitTransaction(ctx context.Context) error {
	tx, ok := ctx.Value(transactionKey).(pgx.Tx)
	if !ok {
		return fmt.Errorf("transaction not found in context")
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (pg *PgContext) RollbackTransaction(ctx context.Context) error {
	tx, ok := ctx.Value(transactionKey).(pgx.Tx)
	if !ok {
		return fmt.Errorf("transaction not found in context")
	}

	if tx == nil {
		return fmt.Errorf("transaction is already nil, can't rollback")
	}

	err := tx.Rollback(ctx)
	if err != nil {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}

	return nil
}

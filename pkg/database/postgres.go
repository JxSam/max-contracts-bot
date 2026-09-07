package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/JxSam/max-contracts-bot/internal/config"
)

type key string

const transactionKey key = "transaction"

type Interface interface {
	TxOrDb(ctx context.Context) dbExecutor
	WithTimeout(timeout time.Duration) (context.Context, context.CancelFunc)
	DefaultTimeoutCtx() (context.Context, context.CancelFunc)
}

type dbExecutor interface {
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
}

type PgContext struct {
	Pool        *pgxpool.Pool
	connTimeout time.Duration
}

func New(dbCfg *config.DatabaseConfig) (*PgContext, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		dbCfg.Host, dbCfg.Port, dbCfg.User, dbCfg.Password, dbCfg.Name)

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}

	return &PgContext{
		Pool:        pool,
		connTimeout: 30 * time.Second,
	}, nil
}

func (pg *PgContext) TxOrDb(ctx context.Context) dbExecutor {
	tx, ok := ctx.Value(transactionKey).(pgx.Tx)
	if !ok {
		return pg.Pool
	}

	return tx
}

func (pg *PgContext) WithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

func (pg *PgContext) DefaultTimeoutCtx() (context.Context, context.CancelFunc) {
	return pg.WithTimeout(pg.connTimeout)
}

func (p *PgContext) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}

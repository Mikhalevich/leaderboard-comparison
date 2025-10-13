package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/Mikhalevich/leaderboard-comparison/internal/domain/scoregenerator"
)

var (
	_ scoregenerator.Repository = (*Postgres)(nil)
)

type Postgres struct {
	db *pgx.Conn
}

func New(db *pgx.Conn) *Postgres {
	return &Postgres{
		db: db,
	}
}

type transactionFn func(tx pgx.Tx) error

func (p *Postgres) transaction(ctx context.Context, txFn transactionFn) error {
	trx, err := p.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("tx begin: %w", err)
	}

	//nolint:errcheck
	defer trx.Rollback(ctx)

	if err := txFn(trx); err != nil {
		return fmt.Errorf("tx fn: %w", err)
	}

	if err := trx.Commit(ctx); err != nil {
		return fmt.Errorf("tx commit: %w", err)
	}

	return nil
}

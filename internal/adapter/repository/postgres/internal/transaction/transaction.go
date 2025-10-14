package transaction

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TrxFn func(tx pgx.Tx) error

func Transaction(ctx context.Context, db *pgxpool.Pool, txFn TrxFn) error {
	trx, err := db.Begin(ctx)
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

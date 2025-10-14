package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/Mikhalevich/leaderboard-comparison/internal/adapter/repository/postgres/internal/transaction"
	"github.com/Mikhalevich/leaderboard-comparison/internal/domain/scoregenerator"
)

func (p *Postgres) AddUserScoreBatch(ctx context.Context, scores []scoregenerator.UserScore) error {
	var (
		query = `
			INSERT INTO 
			score (
				user_id,
				score,
				created_at
			) VALUES (
				$1,
				$2,
				$3
			)
		`
	)

	if err := transaction.Transaction(ctx, p.db, func(tx pgx.Tx) error {
		for _, score := range scores {
			if _, err := tx.Exec(
				ctx,
				query,
				score.UserID,
				score.Score,
				score.CreatedAt,
			); err != nil {
				return fmt.Errorf("insert user score in batch: %w", err)
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction fn: %w", err)
	}

	return nil
}

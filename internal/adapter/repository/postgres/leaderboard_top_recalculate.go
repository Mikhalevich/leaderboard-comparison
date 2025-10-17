package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/Mikhalevich/leaderboard-comparison/internal/domain/leaderboardrecalculator"
)

func (p *Postgres) LeaderboardTopRecalculate(
	ctx context.Context,
	limit int,
) ([]leaderboardrecalculator.LeaderboardEntry, error) {
	var (
		query       = makeLeaderboardRecalculatorTopQueryWithLimit(limit)
		queryRowsFn = func() (pgx.Rows, error) {
			if limit > 0 {
				return p.db.Query(ctx, query, limit)
			}

			return p.db.Query(ctx, query)
		}
	)

	rows, err := queryRowsFn()
	if err != nil {
		return nil, fmt.Errorf("query rows: %w", err)
	}

	defer rows.Close()

	var positions []leaderboardrecalculator.LeaderboardEntry

	for rows.Next() {
		var position leaderboardrecalculator.LeaderboardEntry
		if err := rows.Scan(&position.UserID, &position.Score, &position.Position); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}

		positions = append(positions, position)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return positions, nil
}

func makeLeaderboardRecalculatorTopQueryWithLimit(limit int) string {
	var (
		query = `
			SELECT
				user_id,
				SUM(score) AS user_score,
				ROW_NUMBER() OVER (ORDER BY SUM(score) DESC, MIN(created_at)) AS position
			FROM
				score
			GROUP BY
				user_id
			ORDER BY
				position
			%s
		`
	)

	if limit > 0 {
		return fmt.Sprintf(query, "LIMIT $1")
	}

	return fmt.Sprintf(query, "")
}

package postgres

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/leaderboard-comparison/internal/domain/leaderboard"
)

const (
	halfRankRange = 2
)

func (p *Postgres) LeaderboardByUserID(
	ctx context.Context,
	userID int64,
	limit int,
) ([]leaderboard.LeaderbordEntry, error) {
	var (
		query = `
			WITH cte AS (
				SELECT
					user_id,
					SUM(score) AS user_score,
					ROW_NUMBER() OVER (ORDER BY SUM(score) DESC, MIN(created_at)) AS position
				FROM
					score
				GROUP by
					user_id
				ORDER by
					position
			)
			SELECT
				user_id,
				user_score,
				position
			FROM
				cte
			WHERE
				position >= (
					SELECT
						position - $1
					FROM
						cte
					WHERE
						user_id = $2
				)
			limit $3
		`
	)

	rows, err := p.db.Query(ctx, query, limit/halfRankRange, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("query rows: %w", err)
	}

	defer rows.Close()

	positions := make([]leaderboard.LeaderbordEntry, 0, limit)

	for rows.Next() {
		var position leaderboard.LeaderbordEntry
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

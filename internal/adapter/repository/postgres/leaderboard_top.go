package postgres

import (
	"context"
	"fmt"

	"github.com/Mikhalevich/leaderboard-comparison/internal/domain/leaderboard"
)

func (p *Postgres) LeaderboardTop(ctx context.Context, limit int) ([]leaderboard.LeaderbordEntry, error) {
	var (
		query = `
			SELECT
				user_id,
				SUM(score) AS user_score,
				ROW_NUMBER() OVER (ORDER BY SUM(score) DESC) AS position
			FROM
				score
			GROUP BY
				user_id
			ORDER BY
				position
			LIMIT
				$1
		`
	)

	rows, err := p.db.Query(ctx, query, limit)
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

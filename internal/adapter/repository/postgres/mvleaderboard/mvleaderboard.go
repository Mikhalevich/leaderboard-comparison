package mvleaderboard

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Mikhalevich/leaderboard-comparison/internal/domain/leaderboard"
	"github.com/Mikhalevich/leaderboard-comparison/internal/domain/leaderboardrecalculator"
)

var (
	_ leaderboardrecalculator.LatestTopRepository        = (*MVLeaderboard)(nil)
	_ leaderboardrecalculator.StoreLeaderboardRepository = (*MVLeaderboard)(nil)
	_ leaderboard.LeaderboardRepo                        = (*MVLeaderboard)(nil)
)

const (
	halfRankRange = 2
)

type MVLeaderboard struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *MVLeaderboard {
	return &MVLeaderboard{
		db: db,
	}
}

func (m *MVLeaderboard) LatestTop(ctx context.Context, limit int) ([]leaderboardrecalculator.LeaderboardEntry, error) {
	if _, err := m.db.Exec(ctx, "REFRESH MATERIALIZED VIEW score_leaderboard"); err != nil {
		return nil, fmt.Errorf("refresh mat view: %w", err)
	}

	return nil, nil
}

func (m *MVLeaderboard) StoreLeaderbord(ctx context.Context, top []leaderboardrecalculator.LeaderboardEntry) error {
	return nil
}

func (m *MVLeaderboard) LeaderboardTop(
	ctx context.Context,
	limit int,
) ([]leaderboard.LeaderbordEntry, error) {
	var (
		query = `
			SELECT
				user_id,
				user_score,
				position
			FROM
				score_leaderboard
			limit
				$1
		`
	)

	rows, err := m.db.Query(ctx, query, limit)
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

func (m *MVLeaderboard) LeaderboardByUserID(
	ctx context.Context,
	userID int64,
	limit int,
) ([]leaderboard.LeaderbordEntry, error) {
	var (
		query = `
			SELECT
				user_id,
				user_score,
				position
			FROM
				score_leaderboard
			WHERE
				position >= (
					SELECT
						position - $1
					FROM
						score_leaderboard
					WHERE
						user_id = $2
			)
			limit $3;
		`
	)

	rows, err := m.db.Query(ctx, query, limit/halfRankRange, userID, limit)
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

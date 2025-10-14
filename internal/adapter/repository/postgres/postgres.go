package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Mikhalevich/leaderboard-comparison/internal/domain/leaderboard"
	"github.com/Mikhalevich/leaderboard-comparison/internal/domain/scoregenerator"
)

var (
	_ scoregenerator.Repository   = (*Postgres)(nil)
	_ leaderboard.LeaderboardRepo = (*Postgres)(nil)
)

type Postgres struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Postgres {
	return &Postgres{
		db: db,
	}
}

package postgres

import (
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

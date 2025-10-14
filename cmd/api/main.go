package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Mikhalevich/leaderboard-comparison/internal/adapter/repository/postgres"
	"github.com/Mikhalevich/leaderboard-comparison/internal/adapter/repository/postgres/mvleaderboard"
	"github.com/Mikhalevich/leaderboard-comparison/internal/app/httpapi"
	"github.com/Mikhalevich/leaderboard-comparison/internal/domain/leaderboard"
	"github.com/Mikhalevich/leaderboard-comparison/internal/domain/scoregenerator"
	"github.com/Mikhalevich/leaderboard-comparison/internal/infra"
)

type Config struct {
	Postgres               Postgres `yaml:"postgres" required:"true"`
	IsMVLeaderboardEnabled bool     `yaml:"is_mv_leaderboard_enabled"`
}

type Postgres struct {
	Connection string `yaml:"connection" required:"true"`
}

func main() {
	infra.SetupLogger()

	var cfg Config
	if err := infra.LoadConfig(&cfg); err != nil {
		slog.Error("failed to load config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if err := infra.RunSignalInterruptionFunc(func(ctx context.Context) error {
		slog.Info("api service starting")

		pgxpool, cleanup, err := infra.MakePostgres(ctx, cfg.Postgres.Connection)
		if err != nil {
			return fmt.Errorf("make postgres db: %w", err)
		}

		pgDB := postgres.New(pgxpool)

		defer cleanup()

		if err := httpapi.Start(
			ctx,
			scoregenerator.New(pgDB),
			makeLeaderboard(cfg.IsMVLeaderboardEnabled, pgDB, pgxpool),
		); err != nil {
			return fmt.Errorf("start http api: %w", err)
		}

		slog.Info("api service stopped")

		return nil
	}); err != nil {
		slog.Error("failed run service", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func makeLeaderboard(isMVEnabled bool, pgDB *postgres.Postgres, pgxpool *pgxpool.Pool) *leaderboard.Leaderboard {
	if isMVEnabled {
		return leaderboard.New(mvleaderboard.New(pgxpool))
	}

	return leaderboard.New(pgDB)
}

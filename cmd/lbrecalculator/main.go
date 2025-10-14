package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Mikhalevich/leaderboard-comparison/internal/adapter/repository/postgres/mvleaderboard"
	"github.com/Mikhalevich/leaderboard-comparison/internal/domain/leaderboardrecalculator"
	"github.com/Mikhalevich/leaderboard-comparison/internal/infra"
	"github.com/Mikhalevich/leaderboard-comparison/internal/infra/scheduler"
)

const (
	recalculateAll = -1
)

type Config struct {
	Postgres Postgres      `yaml:"postgres" required:"true"`
	Interval time.Duration `yaml:"interval" required:"true"`
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
		slog.Info("lbrecalculator service starting")

		pgxpool, cleanup, err := infra.MakePostgres(ctx, cfg.Postgres.Connection)
		if err != nil {
			return fmt.Errorf("make postgres db: %w", err)
		}

		defer cleanup()

		var (
			mvPg         = mvleaderboard.New(pgxpool)
			recalculator = leaderboardrecalculator.New(mvPg, mvPg)
		)

		scheduler.Run(
			ctx,
			"leaderboard_relactulator",
			cfg.Interval,
			true,
			func(ctx context.Context) error {
				if err := recalculator.Recalculate(ctx, recalculateAll); err != nil {
					return fmt.Errorf("recalculate: %w", err)
				}

				return nil
			})

		slog.Info("lbrecalculator service stopped")

		return nil
	}); err != nil {
		slog.Error("failed run service", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

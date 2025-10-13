package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/Mikhalevich/leaderboard-comparison/internal/app/httpapi"
	"github.com/Mikhalevich/leaderboard-comparison/internal/domain/scoregenerator"
	"github.com/Mikhalevich/leaderboard-comparison/internal/infra"
)

type Config struct {
	Postgres Postgres `yaml:"postgres" required:"true"`
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

		pgDB, cleanup, err := infra.MakePostgres(ctx, cfg.Postgres.Connection)
		if err != nil {
			return fmt.Errorf("make postgres db: %w", err)
		}

		defer cleanup()

		if err := httpapi.Start(
			ctx,
			scoregenerator.New(pgDB),
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

package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Mikhalevich/leaderboard-comparison/internal/adapter/leaderboardstorer"
	"github.com/Mikhalevich/leaderboard-comparison/internal/adapter/repository/postgres"
	"github.com/Mikhalevich/leaderboard-comparison/internal/domain/leaderboardrecalculator"
	"github.com/Mikhalevich/leaderboard-comparison/internal/infra"
	"github.com/Mikhalevich/leaderboard-comparison/internal/infra/scheduler"
)

const (
	recalculateAll = -1
)

type Config struct {
	Postgres Postgres      `yaml:"postgres" required:"true"`
	Redis    Redis         `yaml:"redis" required:"true"`
	Interval time.Duration `yaml:"interval" required:"true"`
}

type Postgres struct {
	Connection string `yaml:"connection" required:"true"`
}

type Redis struct {
	Addr string `yaml:"addr" required:"true"`
	Pwd  string `yaml:"pwd" required:"true"`
	DB   int    `yaml:"db" required:"true"`
}

func main() {
	infra.SetupLogger()

	var cfg Config
	if err := infra.LoadConfig(&cfg); err != nil {
		slog.Error("failed to load config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if err := infra.RunSignalInterruptionFunc(func(ctx context.Context) error {
		slog.Info("redis_lb_recalculator service starting")

		pgxpool, pgCleanup, err := infra.MakePostgres(ctx, cfg.Postgres.Connection)
		if err != nil {
			return fmt.Errorf("connect to postgres db: %w", err)
		}

		defer pgCleanup()

		rdb, redisCleanup, err := infra.MakeRedis(ctx, cfg.Redis.Addr, cfg.Redis.Pwd, cfg.Redis.DB)
		if err != nil {
			return fmt.Errorf("connect to redis db: %w", err)
		}

		defer redisCleanup()

		var (
			pgDB         = postgres.New(pgxpool)
			redisDB      = leaderboardstorer.New(rdb)
			recalculator = leaderboardrecalculator.New(pgDB, redisDB)
		)

		scheduler.Run(
			ctx,
			"redis_leaderboard_relactulator",
			cfg.Interval,
			true,
			func(ctx context.Context) error {
				if err := recalculator.Recalculate(ctx, recalculateAll); err != nil {
					return fmt.Errorf("recalculate: %w", err)
				}

				return nil
			})

		slog.Info("redis_lb_recalculator service stopped")

		return nil
	}); err != nil {
		slog.Error("failed run service", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

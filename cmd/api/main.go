package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/Mikhalevich/leaderboard-comparison/cmd/api/httpapi"
	"github.com/Mikhalevich/leaderboard-comparison/internal/adapter/leaderboardstorer"
	"github.com/Mikhalevich/leaderboard-comparison/internal/adapter/repository/postgres"
	"github.com/Mikhalevich/leaderboard-comparison/internal/adapter/repository/postgres/mvleaderboard"
	"github.com/Mikhalevich/leaderboard-comparison/internal/domain/leaderboard"
	"github.com/Mikhalevich/leaderboard-comparison/internal/domain/scoregenerator"
	"github.com/Mikhalevich/leaderboard-comparison/internal/infra"
)

type Config struct {
	Postgres                  Postgres `yaml:"postgres" required:"true"`
	Redis                     Redis    `yaml:"redis" required:"false"`
	IsRedisLeaderboardEnabled bool     `yaml:"is_redis_leaderboard_enabled"`
	IsMVLeaderboardEnabled    bool     `yaml:"is_mv_leaderboard_enabled"`
}

type Postgres struct {
	Connection string `yaml:"connection" required:"true"`
}

type Redis struct {
	Addr string `yaml:"addr" required:"false"`
	Pwd  string `yaml:"pwd" required:"false"`
	DB   int    `yaml:"db" required:"false"`
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

		pgxpool, pgCleanup, err := infra.MakePostgres(ctx, cfg.Postgres.Connection)
		if err != nil {
			return fmt.Errorf("connect to postgres db: %w", err)
		}

		defer pgCleanup()

		rdb, redisCleanup, err := makeRedis(ctx, cfg.Redis)
		if err != nil {
			return fmt.Errorf("connect to redis db: %w", err)
		}

		defer redisCleanup()

		pgDB := postgres.New(pgxpool)

		if err := httpapi.Start(
			ctx,
			scoregenerator.New(pgDB),
			makeLeaderboard(cfg.IsRedisLeaderboardEnabled, cfg.IsMVLeaderboardEnabled, pgDB, rdb, pgxpool),
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

func makeLeaderboard(
	isRedisEnabled bool,
	isMVEnabled bool,
	pgDB *postgres.Postgres,
	rdb *redis.Client,
	pgxpool *pgxpool.Pool,
) *leaderboard.Leaderboard {
	if isRedisEnabled {
		return leaderboard.New(leaderboardstorer.New(rdb))
	}

	if isMVEnabled {
		return leaderboard.New(mvleaderboard.New(pgxpool))
	}

	return leaderboard.New(pgDB)
}

func makeRedis(ctx context.Context, cfg Redis) (*redis.Client, func(), error) {
	if cfg.Addr == "" {
		return nil, func() {}, nil
	}

	rdb, cleanup, err := infra.MakeRedis(ctx, cfg.Addr, cfg.Pwd, cfg.DB)
	if err != nil {
		return nil, nil, fmt.Errorf("connect to redis db: %w", err)
	}

	return rdb, cleanup, nil
}

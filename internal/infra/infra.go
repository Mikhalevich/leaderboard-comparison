package infra

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jinzhu/configor"
)

func LoadConfig(cfg any) error {
	configFile := flag.String("config", "config/config.yaml", "config file")
	flag.Parse()

	if err := configor.Load(cfg, *configFile); err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	return nil
}

func SetupLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}

func RunSignalInterruptionFunc(fn func(ctx context.Context) error) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := fn(ctx); err != nil {
		return fmt.Errorf("run fn: %w", err)
	}

	return nil
}

func MakePostgres(ctx context.Context, connection string) (*pgxpool.Pool, func(), error) {
	if connection == "" {
		return nil, func() {}, nil
	}

	conn, err := pgxpool.New(ctx, connection)
	if err != nil {
		return nil, func() {}, fmt.Errorf("pgx connect: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, nil, fmt.Errorf("ping: %w", err)
	}

	return conn, func() {
		conn.Close()
	}, nil
}

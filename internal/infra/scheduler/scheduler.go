package scheduler

import (
	"context"
	"log/slog"
	"time"
)

func Run(
	ctx context.Context,
	workerName string,
	interval time.Duration,
	immediate bool,
	runFn func(ctx context.Context) error,
) {
	if immediate {
		executeFn(ctx, workerName, runFn)
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			executeFn(ctx, workerName, runFn)

		case <-ctx.Done():
			return
		}
	}
}

func executeFn(ctx context.Context, workerName string, fn func(ctx context.Context) error) {
	if err := fn(ctx); err != nil {
		slog.Warn(
			"execute fn error",
			slog.String("worker_name", workerName),
			slog.String("error", err.Error()),
		)
	}
}

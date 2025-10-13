package httpapi

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Mikhalevich/leaderboard-comparison/internal/app/httpapi/internal/handler"
)

const (
	readTimeout      = time.Second * 10
	writeTimeout     = time.Second * 10
	shoutdownTimeout = time.Second * 30
)

func Start(
	ctx context.Context,
	scoreGenerator handler.ScoreGenerator,
) error {
	var (
		mux   = http.NewServeMux()
		hndlr = handler.New(scoreGenerator)
	)

	registerRoutes(mux, hndlr)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			slog.Error("service listen and serve error", slog.String("error", err.Error()))
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shoutdownTimeout)
	defer cancel()

	//nolint:contextcheck
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	return nil
}

func registerRoutes(mux *http.ServeMux, h *handler.Handler) {
	mux.HandleFunc("/generate_test_data", h.GenerateTestData)
}

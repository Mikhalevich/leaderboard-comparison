package leaderboardrecalculator

import (
	"context"
	"fmt"
)

type LeaderboardEntry struct {
	UserID   int64
	Score    int
	Position int
}

type LeaderboardTopRecalculator interface {
	LeaderboardTopRecalculate(ctx context.Context, limit int) ([]LeaderboardEntry, error)
}

type LeaderboardStorer interface {
	StoreLeaderbord(ctx context.Context, top []LeaderboardEntry) error
}

type LeaderboardRecalculator struct {
	latestTopRepo        LeaderboardTopRecalculator
	storeLeaderboardRepo LeaderboardStorer
}

func New(
	latestTopRepo LeaderboardTopRecalculator,
	storeLeaderboardRepo LeaderboardStorer,
) *LeaderboardRecalculator {
	return &LeaderboardRecalculator{
		latestTopRepo:        latestTopRepo,
		storeLeaderboardRepo: storeLeaderboardRepo,
	}
}

func (l *LeaderboardRecalculator) Recalculate(ctx context.Context, limit int) error {
	top, err := l.latestTopRepo.LeaderboardTopRecalculate(ctx, limit)
	if err != nil {
		return fmt.Errorf("latest top repo: %w", err)
	}

	if err := l.storeLeaderboardRepo.StoreLeaderbord(ctx, top); err != nil {
		return fmt.Errorf("store leaderboard repo: %w", err)
	}

	return nil
}

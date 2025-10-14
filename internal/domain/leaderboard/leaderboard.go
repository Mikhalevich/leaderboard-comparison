package leaderboard

import (
	"context"
	"fmt"
)

type LeaderbordEntry struct {
	UserID   int64
	Score    int
	Position int
}

type TopNCalculator interface {
	Top(ctx context.Context, limit int) ([]LeaderbordEntry, error)
}

type Leaderboard struct {
	repo TopNCalculator
}

func New(repo TopNCalculator) *Leaderboard {
	return &Leaderboard{
		repo: repo,
	}
}

func (l *Leaderboard) Top(ctx context.Context, limit int) ([]LeaderbordEntry, error) {
	entries, err := l.repo.Top(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("repo top: %w", err)
	}

	return entries, nil
}

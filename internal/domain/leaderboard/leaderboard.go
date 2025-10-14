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

type LeaderboardRepo interface {
	LeaderboardTop(ctx context.Context, limit int) ([]LeaderbordEntry, error)
	LeaderboardByUserID(ctx context.Context, userID int64, limit int) ([]LeaderbordEntry, error)
}

type Leaderboard struct {
	repo LeaderboardRepo
}

func New(repo LeaderboardRepo) *Leaderboard {
	return &Leaderboard{
		repo: repo,
	}
}

func (l *Leaderboard) Top(ctx context.Context, limit int) ([]LeaderbordEntry, error) {
	entries, err := l.repo.LeaderboardTop(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("repo top: %w", err)
	}

	return entries, nil
}

func (l *Leaderboard) PositionsByUserID(ctx context.Context, userID int64, limit int) ([]LeaderbordEntry, error) {
	entries, err := l.repo.LeaderboardByUserID(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("repo by user id: %w", err)
	}

	return entries, nil
}

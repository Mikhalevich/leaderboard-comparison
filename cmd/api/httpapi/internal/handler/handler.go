package handler

import (
	"context"

	"github.com/Mikhalevich/leaderboard-comparison/internal/domain/leaderboard"
)

type ScoreGenerator interface {
	Generate(ctx context.Context, count int) error
}

type LeaderboardProcessor interface {
	Top(ctx context.Context, limit int) ([]leaderboard.LeaderbordEntry, error)
	ByUserID(ctx context.Context, userID int64, limit int) ([]leaderboard.LeaderbordEntry, error)
}

type Handler struct {
	scoreGenerator       ScoreGenerator
	leaderboardProcessor LeaderboardProcessor
}

func New(
	scoreGenerator ScoreGenerator,
	leaderboardProcessor LeaderboardProcessor,
) *Handler {
	return &Handler{
		scoreGenerator:       scoreGenerator,
		leaderboardProcessor: leaderboardProcessor,
	}
}

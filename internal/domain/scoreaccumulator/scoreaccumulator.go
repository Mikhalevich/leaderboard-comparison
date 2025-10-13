package scoreaccumulator

import (
	"context"
	"fmt"
)

type UserScore struct {
	UserID int64
	Score  int
}

type Repository interface {
	AddUserScore(ctx context.Context, score UserScore) error
}

type ScoreAccumulator struct {
	repo Repository
}

func New(repo Repository) *ScoreAccumulator {
	return &ScoreAccumulator{
		repo: repo,
	}
}

func (s *ScoreAccumulator) Add(ctx context.Context, score UserScore) error {
	if err := s.repo.AddUserScore(ctx, score); err != nil {
		return fmt.Errorf("repo add user score: %w", err)
	}

	return nil
}

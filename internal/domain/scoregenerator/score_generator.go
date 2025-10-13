package scoregenerator

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

const (
	usersRange  = 100000
	scoreRange  = 100
	timeRageSec = 60 * 60 * 24 * 365 * 5
)

type UserScore struct {
	UserID    int64
	Score     int
	CreatedAt time.Time
}

type Repository interface {
	AddUserScoreBatch(ctx context.Context, scores []UserScore) error
}

type ScoreGenerator struct {
	repo Repository
}

func New(repo Repository) *ScoreGenerator {
	return &ScoreGenerator{
		repo: repo,
	}
}

func (s *ScoreGenerator) Generate(ctx context.Context, count int) error {
	scores := generateScores(count)

	if err := s.repo.AddUserScoreBatch(ctx, scores); err != nil {
		return fmt.Errorf("add user score batch: %w", err)
	}

	return nil
}

func generateScores(count int) []UserScore {
	var (
		scores = make([]UserScore, 0, count)
		now    = time.Now()
	)

	for range count {
		scores = append(scores, UserScore{
			//nolint:gosec
			UserID: rand.Int63n(usersRange) + 1,
			//nolint:gosec
			Score: rand.Intn(scoreRange) + 1,
			//nolint:gosec
			CreatedAt: now.Add(-time.Second * time.Duration(rand.Intn(timeRageSec))),
		})
	}

	return scores
}

package handler

import (
	"context"
)

type ScoreGenerator interface {
	Generate(ctx context.Context, count int) error
}

type Handler struct {
	scoreGenerator ScoreGenerator
}

func New(scoreGenerator ScoreGenerator) *Handler {
	return &Handler{
		scoreGenerator: scoreGenerator,
	}
}

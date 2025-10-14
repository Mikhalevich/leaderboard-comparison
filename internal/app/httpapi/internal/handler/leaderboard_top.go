package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Mikhalevich/leaderboard-comparison/internal/domain/leaderboard"
)

const (
	leaderboardLimit = 10
)

func (h *Handler) LeaderboardTop(rspWriter http.ResponseWriter, req *http.Request) {
	positions, err := h.leaderboardProcessor.Top(req.Context(), leaderboardLimit)
	if err != nil {
		slog.Error("leaderboard top", slog.String("error", err.Error()))
		http.Error(rspWriter, "leaderboard top error", http.StatusInternalServerError)

		return
	}

	rspWriter.Header().Add("Content-Type", "application/json")

	jsonPositions := convertToLeaderboardEntryJSONPayload(positions)

	if err := json.NewEncoder(rspWriter).Encode(&jsonPositions); err != nil {
		slog.Error("leaderboard json encode", slog.String("error", err.Error()))
		http.Error(rspWriter, "leaderboard json encode error", http.StatusInternalServerError)

		return
	}
}

type LeaderboardEntryJSONPayload struct {
	UserID   int64 `json:"user_id"`
	Score    int   `json:"score"`
	Position int   `json:"position"`
}

func convertToLeaderboardEntryJSONPayload(dbPositions []leaderboard.LeaderbordEntry) []LeaderboardEntryJSONPayload {
	jsonPositions := make([]LeaderboardEntryJSONPayload, 0, len(dbPositions))

	for _, v := range dbPositions {
		jsonPositions = append(jsonPositions, LeaderboardEntryJSONPayload{
			UserID:   v.UserID,
			Score:    v.Score,
			Position: v.Position,
		})
	}

	return jsonPositions
}

package handler

import (
	"encoding/json"
	"log/slog"
	"math/rand"
	"net/http"
)

const (
	userIDRange = 100000
)

func (h *Handler) LeaderboardByUserIDRandom(rspWriter http.ResponseWriter, req *http.Request) {
	//nolint:gosec
	userID := rand.Int63n(userIDRange) + 1

	positions, err := h.leaderboardProcessor.ByUserID(req.Context(), userID, leaderboardLimit)
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

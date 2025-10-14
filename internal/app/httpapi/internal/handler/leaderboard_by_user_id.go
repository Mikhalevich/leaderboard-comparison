package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

func (h *Handler) LeaderboardByUserID(rspWriter http.ResponseWriter, req *http.Request) {
	userID, err := userIDFromRequest(req)
	if err != nil {
		slog.Error("invalid user id", slog.String("error", err.Error()))
		http.Error(rspWriter, "invalid user id", http.StatusBadRequest)

		return
	}

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

func userIDFromRequest(r *http.Request) (int64, error) {
	rawUserID := r.PathValue("user_id")
	if rawUserID == "" {
		return 0, errors.New("no user id found")
	}

	userID, err := strconv.ParseInt(rawUserID, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("user_id convert: %w", err)
	}

	return userID, nil
}

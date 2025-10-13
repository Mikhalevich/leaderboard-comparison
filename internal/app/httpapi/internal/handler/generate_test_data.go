package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type GenerateTestDataRequest struct {
	RowsCount int `json:"rows_count"`
}

func (h *Handler) GenerateTestData(w http.ResponseWriter, r *http.Request) {
	var req GenerateTestDataRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("json decode", slog.String("error", err.Error()))
		http.Error(w, "json decode erorr", http.StatusBadRequest)

		return
	}

	if req.RowsCount <= 0 {
		http.Error(w, "invalid rows count", http.StatusBadRequest)

		return
	}

	if err := h.scoreGenerator.Generate(r.Context(), req.RowsCount); err != nil {
		slog.Error("score generate error", slog.String("error", err.Error()))
		http.Error(w, "generate data error", http.StatusInternalServerError)

		return
	}
}

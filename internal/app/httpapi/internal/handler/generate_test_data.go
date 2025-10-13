package handler

import (
	"log/slog"
	"net/http"
)

const (
	defaultTestDataCount = 10000
)

func (h *Handler) GenerateTestData(w http.ResponseWriter, r *http.Request) {
	if err := h.scoreGenerator.Generate(r.Context(), defaultTestDataCount); err != nil {
		slog.Error("score generate error", slog.String("error", err.Error()))
		http.Error(w, "generate data error", http.StatusInternalServerError)

		return
	}
}

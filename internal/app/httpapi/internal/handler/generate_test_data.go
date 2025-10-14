package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type GenerateTestDataRequest struct {
	RowsCount int `json:"rows_count"`
}

func (h *Handler) GenerateTestData(rspWriter http.ResponseWriter, httpReq *http.Request) {
	var req GenerateTestDataRequest
	if err := json.NewDecoder(httpReq.Body).Decode(&req); err != nil {
		slog.Error("json decode", slog.String("error", err.Error()))
		http.Error(rspWriter, "json decode erorr", http.StatusBadRequest)

		return
	}

	if req.RowsCount <= 0 {
		http.Error(rspWriter, "invalid rows count", http.StatusBadRequest)

		return
	}

	if err := h.scoreGenerator.Generate(httpReq.Context(), req.RowsCount); err != nil {
		slog.Error("score generate error", slog.String("error", err.Error()))
		http.Error(rspWriter, "generate data error", http.StatusInternalServerError)

		return
	}
}

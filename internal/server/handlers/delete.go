package handlers

import (
	"context"
	"httpfs/internal/entities"
	"httpfs/internal/logger"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Deleter интерфейс удаления файла
type Deleter interface {
	// удаление файла
	Delete(context.Context, entities.Hash) error
}

func HandleDelete(d Deleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := chi.URLParam(r, "hash")
		if hash == "" {
			respondJson(w, http.StatusBadRequest, nil)
			return
		}

		ctx := logger.LogWithKvContext(r.Context(), "hash", hash)

		if err := d.Delete(ctx, entities.Hash(hash)); err != nil {
			slog.ErrorContext(ctx,
				"fault delete",
				"err", err,
			)

			respondJson(w, http.StatusInternalServerError, errorMessage{
				Message: "fault delete file",
			})
			return
		}
	}
}

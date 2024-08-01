package handlers

import (
	"context"
	"errors"
	"httpfs/internal/entities"
	"httpfs/internal/logger"
	"httpfs/internal/storage"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Downloader интерфейс загрузки файла
type Downloader interface {
	// Download загрузка файла
	Download(context.Context, entities.Hash, io.Writer) error
}

// HandleDownload обработчик загрузки файла
func HandleDownload(d Downloader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		has := chi.URLParam(r, "hash")
		if has == "" {
			respondJson(w, http.StatusBadRequest, errorMessage{
				Message: "invalid hash",
			})
			return
		}

		ctx := logger.LogWithKvContext(r.Context(), "hash", has)

		if err := d.Download(ctx, entities.Hash(has), w); err != nil {

			if errors.Is(err, storage.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			slog.ErrorContext(ctx,
				"fault download",
				"err", err,
			)

			respondJson(w, http.StatusInternalServerError, errorMessage{
				Message: "fault download file",
			})
			return
		}
	}
}

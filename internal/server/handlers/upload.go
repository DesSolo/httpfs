package handlers

import (
	"bytes"
	"context"
	"httpfs/internal/entities"
	"io"
	"log/slog"
	"net/http"
)

// Uploader интерфейс загрузки файла
type Uploader interface {
	// Upload загрузка файла
	Upload(context.Context, io.Reader) (entities.Hash, error)
}

// CallBack интерфейс обработчика
type CallBack interface {
	// предварительная обработка
	Pre(context.Context, io.Reader) error
	// обработка после успешной загрузки файла
	Post(context.Context, entities.Hash) error
}

// uploadResponse структура ответа
type uploadResponse struct {
	// хэш загруженного файла
	Hash entities.Hash `json:"hash"`
}

// HandleUpload обработчик загрузки файла
func HandleUpload(cb CallBack, u Uploader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if r.Body == nil {
			respondJson(w, http.StatusBadRequest, errorMessage{
				Message: "empty body",
			})
			return
		}

		// тут вычитываем все в память
		// Не самое лучшее решение, если размер загружаемого файла будет больший
		// нужно использовать io.Pipe и io.TeeReader
		data, err := io.ReadAll(r.Body)
		if err != nil {
			slog.ErrorContext(ctx,
				"fault read body",
				"err", err,
			)
			respondJson(w, http.StatusBadRequest, errorMessage{
				Message: "fault read body",
			})
			return
		}

		if err := cb.Pre(ctx, bytes.NewBuffer(data)); err != nil {

			slog.ErrorContext(ctx,
				"fault pre upload",
				"err", err,
			)

			respondJson(w, http.StatusBadRequest, errorMessage{
				Message: err.Error(),
			})
			return
		}

		hash, err := u.Upload(ctx, bytes.NewBuffer(data))
		if err != nil {

			slog.ErrorContext(ctx,
				"fault upload",
				"err", err,
			)

			respondJson(w, http.StatusInternalServerError, errorMessage{
				Message: "fault upload file",
			})
			return
		}

		if err := cb.Post(ctx, hash); err != nil {
			slog.ErrorContext(ctx,
				"fault post upload",
				"err", err,
			)

			respondJson(w, http.StatusInternalServerError, errorMessage{
				Message: "fault post upload file",
			})
			return
		}

		respondJson(w, http.StatusOK, uploadResponse{
			Hash: hash,
		})
	}
}

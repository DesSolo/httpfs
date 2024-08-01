package storage

import (
	"context"
	"httpfs/internal/entities"
	"io"
)

// Storage интерфейс хранилища
type Storage interface {
	// Upload загрузка файла
	Upload(context.Context, io.Reader) (entities.Hash, error)
	// Download скачивание файла по хешу
	Download(context.Context, entities.Hash, io.Writer) error
	// Delete удаление файла по хешу
	Delete(context.Context, entities.Hash) error
}

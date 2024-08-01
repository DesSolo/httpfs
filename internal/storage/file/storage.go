package file

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"httpfs/internal/entities"
	"httpfs/internal/storage"
	"io"
	"log/slog"
	"os"
	"path"
)

// tempPattern имени временного файла
const tempPattern = "httpfs-tmp"

// Storage файловый хранилище
type Storage struct {
	// путь к постоянному хранилищу
	basePath string
	// путь к временному хранилищу
	temporaryPath string
}

// NewStorage конструктор
func NewStorage(basePath, temporaryPath string) *Storage {
	return &Storage{
		basePath:      basePath,
		temporaryPath: temporaryPath,
	}
}

// fullPath возвращает возможный путь к файлу внутри постоянного хранилища
func (s *Storage) fullPath(h entities.Hash) (string, error) {
	if len(h) < 2 {
		return "", fmt.Errorf("hash too short")
	}

	hashString := h.String()

	// путь к файлу внутри постоянного хранилища - <basePath>/<2 первых символа хэша>/<хэш файла>
	return path.Join(s.basePath, hashString[:2], hashString), nil
}

// Upload загрузка файла
func (s *Storage) Upload(ctx context.Context, r io.Reader) (entities.Hash, error) {
	file, err := os.CreateTemp(s.temporaryPath, tempPattern)
	if err != nil {
		return "", fmt.Errorf("create temp file err: %w", err)
	}

	// TODO: можно вынести в интерфейс чтоб была возможность выбора алгоритма хэширования
	h := sha1.New()

	// сразу считаем хэш
	w := io.MultiWriter(file, h)

	if _, err := io.Copy(w, r); err != nil {
		return "", fmt.Errorf("io copy err: %w", err)
	}

	if err := file.Close(); err != nil {
		return "", fmt.Errorf("close file err: %w", err)
	}

	hash := entities.Hash(
		fmt.Sprintf("%x", h.Sum(nil)),
	)

	filePath, err := s.fullPath(hash)
	if err != nil {
		return "", fmt.Errorf("full path err: %w", err)
	}

	if err := os.MkdirAll(path.Dir(filePath), os.ModePerm); err != nil {
		return "", fmt.Errorf("create dir err: %w", err)
	}

	// TODO: добавить проверку если такой файл уже есть в basePath

	// переносим файл в постоянное хранилище
	if err := os.Rename(file.Name(), filePath); err != nil {
		return "", fmt.Errorf("move file err: %w", err)
	}

	return hash, nil
}

// Download загрузка файла
func (s *Storage) Download(ctx context.Context, h entities.Hash, w io.Writer) error {
	filePath, err := s.fullPath(h)
	if err != nil {
		return fmt.Errorf("full path err: %w", err)
	}

	slog.DebugContext(ctx,
		"path",
		"path", filePath,
	)

	file, err := os.Open(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("%w err: %w", storage.ErrNotFound, err)
		}

		return fmt.Errorf("open file err: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(w, file); err != nil {
		return fmt.Errorf("io copy err: %w", err)
	}

	return nil
}

// Delete удаление файла
func (s *Storage) Delete(ctx context.Context, h entities.Hash) error {
	filePath, err := s.fullPath(h)
	if err != nil {
		return fmt.Errorf("full path err: %w", err)
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("stat file err: %w", err)
	}

	if info.IsDir() {
		return fmt.Errorf("not a file")
	}

	return os.Remove(filePath)
}

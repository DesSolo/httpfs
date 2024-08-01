package simple

import (
	"context"
	"fmt"
	"httpfs/internal/callback"
	"io"
	"log/slog"
)

// FileSizePre принимает только файлы в диапазоне min..max
// это не лучшая реализация получения размера при копировании в discard
// эффективнее просто не читать больше чем нужно
func FileSizePre(min, max int64) callback.PreCallBackFunc {
	return func(ctx context.Context, r io.Reader) error {
		size, err := io.Copy(io.Discard, r)
		if err != nil {
			return fmt.Errorf("io copy err: %w", err)
		}

		slog.DebugContext(ctx,
			"file size",
			"min", min,
			"max", max,
			"actual", size,
		)

		if size < min || size > max {
			return fmt.Errorf("file size err: %d", size)
		}

		return nil
	}
}

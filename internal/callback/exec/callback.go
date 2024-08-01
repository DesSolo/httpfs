package exec

import (
	"context"
	"fmt"
	"httpfs/internal/callback"
	"httpfs/internal/entities"
	"log/slog"
	"os"
	"os/exec"
)

// envNameHash имя переменной окружения с хешем
const envNameHash = "HASH"

// PostCallBackFunc callback на выполнение команды
func Post(command string, args ...string) callback.PostCallBackFunc {
	return func(ctx context.Context, h entities.Hash) error {

		if err := os.Setenv(envNameHash, h.String()); err != nil {
			return fmt.Errorf("set env err: %w", err)
		}

		slog.DebugContext(ctx,
			"exec",
			"command", command,
			"args", args,
			"hash", h,
		)

		return exec.CommandContext(ctx, command, args...).Run()
	}
}

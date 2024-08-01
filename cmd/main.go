package main

import (
	"fmt"
	"httpfs/config"
	"httpfs/internal/callback"
	"httpfs/internal/callback/exec"
	"httpfs/internal/callback/simple"
	"httpfs/internal/logger"
	"httpfs/internal/server"
	"httpfs/internal/storage"
	"httpfs/internal/storage/file"
	"log/slog"
	"os"
)

// loadConfig загружает конфигурацию
func loadConfig() (*config.Config, error) {
	configFilePath := os.Getenv("CONFIG_FILE_PATH")

	if configFilePath == "" {
		configFilePath = "/etc/httpfs/config.yaml"
	}

	return config.NewConfigFromFile(configFilePath)
}

// configureLogger настраивает логгер
func configureLogger(cfg *config.Config) error {
	o := slog.HandlerOptions{
		AddSource: cfg.Logging.Option.AddSource,
		Level:     slog.Level(cfg.Logging.Option.Level),
	}

	var h slog.Handler

	switch cfg.Logging.Handler {
	case "text":
		h = slog.NewTextHandler(os.Stdout, &o)
	case "json":
		h = slog.NewJSONHandler(os.Stdout, &o)
	default:
		return fmt.Errorf("handler: %q not supported", cfg.Logging.Handler)
	}

	h = logger.NewLogContextHandler(h)

	slog.SetDefault(
		slog.New(h),
	)

	return nil
}

// loadStorage загружает хранилище
func loadStorage(cfg *config.Config) (storage.Storage, error) {
	switch cfg.Storage.Kind {
	case "file":
		return file.NewStorage(cfg.Storage.File.BasePath, cfg.Storage.File.TemporaryPath), nil
	default:
		return nil, fmt.Errorf("storage: %q not supported", cfg.Storage.Kind)
	}
}

func loadCallback(cfg *config.Config) (*callback.CallBack, error) {
	cb := callback.New()

	for _, p := range cfg.Callback.Pre {
		switch p.Kind {
		case "file_size":
			cb.RegisterPre(simple.FileSizePre(p.FileSize.Min, p.FileSize.Max))
		default:
			return nil, fmt.Errorf("callback: %q not supported", p.Kind)
		}
	}

	for _, p := range cfg.Callback.Post {
		switch p.Kind {
		case "exec":
			cb.RegisterPost(exec.Post(p.Exec.Command, p.Exec.Args...))
		default:
			return nil, fmt.Errorf("callback: %q not supported", p.Kind)
		}
	}

	return cb, nil
}

// loadServer загружает HTTP сервер
func loadServer(cb *callback.CallBack, st storage.Storage) *server.Server {
	srv := server.NewServer()
	srv.LoadRoutes(cb, st)
	return srv
}

// fatal выдаёт ошибку и завершает работу
func fatal(message string, err error) {
	slog.Error(message, "err", err)
	os.Exit(1)
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		fatal("fault load config err", err)
	}

	if err := configureLogger(cfg); err != nil {
		fatal("fault configure logger err", err)
	}

	cb, err := loadCallback(cfg)
	if err != nil {
		fatal("fault load callback err", err)
	}

	st, err := loadStorage(cfg)
	if err != nil {
		fatal("fault load storage err", err)
	}

	srv := loadServer(cb, st)

	slog.Info("server running", "addr", cfg.Address)

	if err := srv.Run(cfg.Address); err != nil {
		fatal("fault run server err", err)
	}
}

package logger

import (
	"context"
	"log/slog"

	"github.com/go-chi/chi/v5/middleware"
)

// LogContextHandler обертка для логгера с использованием контекста
type LogContextHandler struct {
	slog.Handler
}

// NewLogContextHandler конструктор
func NewLogContextHandler(h slog.Handler) *LogContextHandler {
	return &LogContextHandler{
		Handler: h,
	}
}

// Handle обработчик slog.Record для проставления значений из context
func (h LogContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if reqID := middleware.GetReqID(ctx); reqID != "" {
		r.Add("reqID", reqID)
	}

	if kv, ok := ctx.Value(logCtxKey).(ctxKV); ok {
		r.Add("kv", kv)
	}

	return h.Handler.Handle(ctx, r)
}

// WithAttrs возвращает обработчик с добавленными аттрибутами (для совместимости с slog.Handler)
func (h LogContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return LogContextHandler{
		Handler: h.Handler.WithAttrs(attrs),
	}
}

// WithGroup возвращает обработчик с добавленной группой (для совместимости с slog.Handler)
func (h LogContextHandler) WithGroup(name string) slog.Handler {
	return LogContextHandler{
		Handler: h.Handler.WithGroup(name),
	}
}

package server

import (
	"httpfs/internal/callback"
	"httpfs/internal/server/handlers"
	"httpfs/internal/storage"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Server HTTP сервер
type Server struct {
	// mux маршрутизатор
	mux chi.Router
}

// NewServer создает новый HTTP сервер
func NewServer() *Server {
	return &Server{
		mux: chi.NewRouter(),
	}
}

// LoadRoutes загружает маршруты
func (s *Server) LoadRoutes(cb *callback.CallBack, st storage.Storage) {
	s.mux.Route("/", func(r chi.Router) {
		r.Use(
			middleware.RealIP,
			middleware.RequestID,
			middleware.Logger,
			middleware.Recoverer,
			// тут можно добавить rate limiter например https://github.com/go-chi/httprate
			// но такое лучше делать на балансировщике
		)

		r.Put("/{file}", handlers.HandleUpload(cb, st))

		r.Get("/{hash}", handlers.HandleDownload(st))
		r.Delete("/{hash}", handlers.HandleDelete(st))

		r.Get("/health", handlers.HandleHealth())
	})
}

// Run запускает сервер на указанном адресе
func (s *Server) Run(addr string) error {
	return http.ListenAndServe(addr, s.mux)
}

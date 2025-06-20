package server

import (
	"context"
	"net/http"
	"silkroad/m/internal/config"
)

type Server struct {
	httServer *http.Server
	config    config.ServerConfig
}

func New(cfg config.ServerConfig) *Server {
	return &Server{
		config: cfg,
	}
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: s.config.MaxHeaderBytes,
		ReadTimeout:    s.config.ReadTimeout,
		WriteTimeout:   s.config.WriteTimeout,
	}

	return s.httServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httServer.Shutdown(ctx)
}

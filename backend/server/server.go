package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mvd-inc/anibliss/internal/config"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.ServerConfig, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			MaxHeaderBytes: 2048,
			ReadTimeout:    time.Duration(1000000000000000), // TODO брать из конфига
			WriteTimeout:   time.Duration(1000000000000000),
			Addr:           fmt.Sprintf("0.0.0.0:%d", cfg.Port),
			Handler:        handler,
		},
	}
}
func (s *Server) Serve() error {
	return s.httpServer.ListenAndServe()
}

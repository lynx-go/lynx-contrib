package http

import (
	"context"
	"github.com/lynx-go/lynx"
	"log/slog"
	"net/http"
)

func New(addr string, handler http.Handler) *Server {
	return &Server{
		addr:    addr,
		handler: handler,
	}
}

type Server struct {
	addr    string
	handler http.Handler
	server  *http.Server
}

func (s *Server) Name() string {
	return "http"
}

func (s *Server) Start(ctx context.Context) error {
	logger := slog.Default()
	go func() {
		s.server = &http.Server{Addr: s.addr, Handler: s.handler}
		logger.Info("http server starting", "addr", s.addr)
		if err := s.server.ListenAndServe(); err != nil {
			logger.Error("start http server failed", "error", err)
		}
	}()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Close()
}

var _ lynx.Server = new(Server)

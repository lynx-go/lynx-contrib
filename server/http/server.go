package http

import (
	"context"
	"github.com/lynx-go/lynx"
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

func (s *Server) Start(ctx context.Context) error {
	s.server = &http.Server{Addr: s.addr, Handler: s.handler}
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) Name() string {
	return "http"
}

var _ lynx.Hook = new(Server)

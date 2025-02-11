package http

import (
	"context"
	"github.com/lynx-go/lynx"
	"github.com/lynx-go/lynx/hook"
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
	*hook.HookBase
}

func (s *Server) OnStart(ctx context.Context) error {
	s.server = &http.Server{Addr: s.addr, Handler: s.handler}
	return s.server.ListenAndServe()
}

func (s *Server) OnStop(ctx context.Context) {
	_ = s.server.Shutdown(ctx)
}

func (s *Server) Name() string {
	return "http"
}

var _ lynx.Server = new(Server)

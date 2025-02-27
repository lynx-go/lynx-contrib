package http

import (
	"context"
	"github.com/lynx-go/lynx/hook"
	"net/http"
	"sync"
)

type Option func(srv *Server)

func WithAddr(addr string) Option {
	return func(srv *Server) {
		srv.addr = addr
	}
}

func WithHandler(h http.Handler) Option {
	return func(srv *Server) {
		srv.handler = h
	}
}

func (s *Server) ensureDefaults() {
	if s.addr == "" {
		s.addr = ":8080"
	}
	if s.handler == nil {
		s.handler = http.DefaultServeMux
	}
}

func New(opts ...Option) *Server {
	s := &Server{}
	for _, opt := range opts {
		opt(s)
	}
	s.ensureDefaults()
	return s
}

type Server struct {
	addr    string
	handler http.Handler
	server  *http.Server

	started bool
	mux     sync.Mutex
}

func (s *Server) Status() (hook.Status, error) {
	if s.started {
		return hook.StatusStarted, nil
	}
	return hook.StatusUnstart, nil
}

func (s *Server) Start(ctx context.Context) error {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.started = true
	s.server = &http.Server{Addr: s.addr, Handler: s.handler}
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.started = false
	return s.server.Shutdown(ctx)
}

func (s *Server) Name() string {
	return "http"
}

var _ hook.Hook = new(Server)

package http

import (
	"context"
	"errors"
	"github.com/lynx-go/lynx/hook"
	"github.com/lynx-go/x/log"
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
	log.InfoContext(ctx, "http server starting", "addr", s.addr)
	if err := s.server.ListenAndServe(); err != nil {
		s.started = false
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.started = false
	return s.server.Shutdown(ctx)
}

func (s *Server) Name() string {
	return "http"
}

var _ hook.Hook = new(Server)

package gin

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lynx-go/lynx/hook"
	"net/http"
	"sync"
)

type Option struct {
	Addr         string
	PanicOnError bool
}

type Server struct {
	o       Option
	srv     *http.Server
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
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.started = false
	return s.srv.Shutdown(ctx)
}

func (s *Server) Name() string {
	return "gin-server"
}

var _ hook.Hook = new(Server)

type MountRoutes func(r *gin.Engine)

func New(o Option, mount MountRoutes) *Server {
	router := gin.Default()
	mount(router)
	srv := &http.Server{
		Addr:    o.Addr,
		Handler: router,
	}
	s := &Server{
		o:   o,
		srv: srv,
	}
	return s
}

package gin

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lynx-go/lynx"
	"github.com/lynx-go/lynx/hook"
	"net/http"
)

type Option struct {
	Addr         string
	PanicOnError bool
}

type Server struct {
	*hook.HookBase
	o   Option
	srv *http.Server
}

func (s *Server) OnStart(ctx context.Context) error {
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) OnStop(ctx context.Context) {
	_ = s.srv.Shutdown(ctx)
}

func (s *Server) IgnoreForCLI() bool {
	return false
}

func (s *Server) Name() string {
	return "gin"
}

var _ lynx.Server = new(Server)

type MountFunc func(r gin.IRoutes)

func New(o Option, mount MountFunc) *Server {
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

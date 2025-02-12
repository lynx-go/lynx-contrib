package gin

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lynx-go/lynx/lifecycle"
	"net/http"
)

type Option struct {
	Addr         string
	PanicOnError bool
}

type Server struct {
	o   Option
	srv *http.Server
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) {
	_ = s.srv.Shutdown(ctx)
}

func (s *Server) IgnoreCLI() bool {
	return true
}

func (s *Server) Name() string {
	return "gin"
}

var _ lifecycle.Service = new(Server)

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

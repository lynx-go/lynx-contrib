package gin

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lynx-go/lynx"
	"log/slog"
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

func (s *Server) Name() string {
	return "gin"
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		if err := s.srv.ListenAndServe(); err != nil {
			slog.Error("start gin server failed", "error", err)
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			if s.o.PanicOnError {
				panic(err)
			}
		}
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
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

package gin

import (
	"context"
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestServer_Start(t *testing.T) {
	s := New(Option{Addr: ":8080"}, func(r gin.IRoutes) {
		r.Any("/hello", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "hello world"})
		})
	})
	ctx := context.TODO()
	if err := s.Start(ctx); err != nil {
		t.Fatal(err)
	}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	s.Stop(ctx)
	time.Sleep(1 * time.Second)
}

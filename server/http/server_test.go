package http

import (
	"github.com/lynx-go/lynx"
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello World"))
	}))
	server := New(":8888", mux)
	app := lynx.New(lynx.WithName("lynx-http"), lynx.WithServer(server))
	app.Run()
}

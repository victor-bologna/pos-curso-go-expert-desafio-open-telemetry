package webserver

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"
)

type WebServer struct {
	Handlers      map[string]http.HandlerFunc
	WebServerPort string
	srv           *http.Server
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Handlers:      make(map[string]http.HandlerFunc),
		WebServerPort: serverPort,
	}
}

func (ws *WebServer) AddHandler(method, path string, handler http.HandlerFunc) {
	ws.Handlers[strings.ToUpper(method)+" "+path] = handler
}

func (ws *WebServer) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	for path, handler := range ws.Handlers {
		mux.HandleFunc(path, handler)
	}
	ws.srv = &http.Server{
		Addr:         ":" + ws.WebServerPort,
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      mux,
	}
	return ws.srv.ListenAndServe()
}

func (ws *WebServer) Shutdown(ctx context.Context) error {
	return ws.srv.Shutdown(ctx)
}

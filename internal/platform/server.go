package platform

import (
	"context"
	"net/http"
)

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer(srv *http.Server) *HTTPServer {
	return &HTTPServer{
		server: srv,
	}
}

func (s *HTTPServer) Start() error {
	return s.server.ListenAndServe()
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.server.Close()
}

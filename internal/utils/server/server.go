package server

import "net/http"

type Server struct {
	server http.Server
}

func New(addr string, handler http.Handler) *Server {
	return &Server{
		server: http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

func (s *Server) Start() {
	s.server.ListenAndServe()
}

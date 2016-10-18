package server

import (
	"net/http"
)

const (
	ServerVersion = 0.1
)

// Server with request multiplexer
type Server struct {
	mux *http.ServeMux
}

// New creates instance of Server
func New() *Server {
	return &Server{mux: http.NewServeMux()}
}

// Start registers available routes and starts server
func (s *Server) Start() error {
	LogInfo("Starting server")
	s.mux.HandleFunc("/", rootHandler)
	s.mux.HandleFunc("/numbers", numbersHandler)
	return http.ListenAndServe(":8080", s.mux)
}

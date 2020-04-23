package server

import (
	"net/http"

	"github.com/go-chi/chi"
	kitlog "github.com/go-kit/kit/log"
)

type Server struct {
	Logger kitlog.Logger
	router chi.Router
}

func New(logger kitlog.Logger) *Server {
	s := &Server{Logger: logger}
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

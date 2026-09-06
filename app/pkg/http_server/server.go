package httpserver

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Logger interface {
	Error(string, error)
}

type Server struct {
	r chi.Router
	l Logger
}

func NewServer(r chi.Router, l Logger) *Server {
	return &Server{r: r, l: l}
}

func (s *Server) Serve(port int) error {
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), s.r)
	if err != nil {
		s.l.Error("server was stopped", err)
		return err
	}

	return nil
}

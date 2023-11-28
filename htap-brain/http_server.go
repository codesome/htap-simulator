package main

import (
	"net/http"
)

type Server struct {
	mux     *http.ServeMux
	writer  *Writer
	querier *Querier
}

func NewServer(w *Writer, q *Querier) *Server {
	return &Server{
		writer:  w,
		querier: q,
	}
}

func (s *Server) ListenAndServe() error {
	s.mux = http.NewServeMux()
	s.mux.HandleFunc("/write", s.write)
	s.mux.HandleFunc("/read", s.read)
	return http.ListenAndServe(":3333", s.mux)
}

func (s *Server) write(w http.ResponseWriter, r *http.Request) {

}
func (s *Server) read(w http.ResponseWriter, r *http.Request) {

}

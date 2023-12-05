package main

import (
	"net/http"
)

type Server struct {
	mux  *http.ServeMux
	htap *HTAPBrain
}

func NewServer(htap *HTAPBrain) *Server {
	return &Server{
		htap: htap,
	}
}

func (s *Server) ListenAndServe() error {
	s.mux = http.NewServeMux()
	s.mux.HandleFunc("/write", s.write)
	s.mux.HandleFunc("/read", s.read)
	return http.ListenAndServe(":3333", s.mux)
}

func (s *Server) write(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(err.Error()))
		return
	}
	query := r.Form.Get("query")
	err = s.htap.Write(query)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(err.Error()))
		return
	}
}
func (s *Server) read(w http.ResponseWriter, r *http.Request) {

}

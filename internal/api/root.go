package api

import (
	"net/http"
)

func (s *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

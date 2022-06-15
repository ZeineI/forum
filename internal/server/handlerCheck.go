package server

import "net/http"

func (s *Server) Check(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" || r.URL.Path == "/" {
		ErrorHandler(w, http.StatusNotFound)
		return
	}
}

package server

import (
	"net/http"
)

//Middleware
func (s *Server) MiddlewareLoginCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rightCookie := s.CheckCookie(w, r)
		if rightCookie {
			next.ServeHTTP(w, r)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})
}

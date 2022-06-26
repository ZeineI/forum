package server

import (
	"log"
	"net/http"
)

//check cookie
func (s *Server) CheckCookie(w http.ResponseWriter, r *http.Request) bool {
	var rightCookie bool
	cookie, err := r.Cookie("session")
	if err != nil {
		log.Println("CheckCookie: There is no cookies for this user")
		return false
	}
	if err := s.cookies.CheckSession(cookie.Value); err != nil {
		log.Println("There is cookie for this user BUT OLD")
		//delete cookie
		cookie.MaxAge = -1
		rightCookie = false
		http.SetCookie(w, cookie)
	} else {
		rightCookie = true
	}

	return rightCookie
}

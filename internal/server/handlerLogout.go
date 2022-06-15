package server

import (
	"log"
	"net/http"
)

func (s *Server) LogOut(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/logout" {
		ErrorHandler(w, http.StatusBadRequest)
		return
	}
	AlreadyLogOut := s.CheckCookie(w, r)
	if !AlreadyLogOut {
		http.Redirect(w, r, "/main", http.StatusSeeOther)
		return
	}
	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
	cookies, err := r.Cookie("session")
	if err != nil {
		log.Println("There is no cookies for this user")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// coockieVal := cookies.Value
	s.cookies.DeleteCookie(cookies.Value)
	cookies.MaxAge = -1
	http.SetCookie(w, cookies)

	//if auth by github
	// place, err := s.Db.GetPlace(coockieVal)
	// if err != nil {
	// 	ErrorHandler(w, http.StatusInternalServerError)
	// 	return
	// }

	// if place == "github" {
	// 	s.Db.DeleteUser(coockieVal)
	// }

	http.Redirect(w, r, "/main", http.StatusSeeOther)
}

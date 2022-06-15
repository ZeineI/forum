package server

import (
	"log"
	"net/http"
	"text/template"
)

func (s *Server) MainHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/main" {
		log.Println("main handler url error")
		ErrorHandler(w, http.StatusBadRequest)
		return
	}
	if r.Method != http.MethodGet {
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
	files := []string{
		"./static/mainPage.html",
		"./static/header.html",
	}
	mainPage, err := template.ParseFiles(files...)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	auth, err := s.AuthPostOperations(w, r)
	if err != nil {
		log.Println("main handler:", err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	if err := mainPage.Execute(w, auth); err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
}

package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func (s *Server) MainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AAAAAAAA")
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
		"./frontend/html/mainPage.html",
		"./frontend/html/header.html",
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

package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/ZeineI/forum/internal/models"
)

func (s *Server) MyPosts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/myPosts" {
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
	myPostsPage, err := template.ParseFiles(files...)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	isAuthStruct := new(models.IsAuthStruct)
	isAuthStruct.IsAuth = true

	cookie, err := r.Cookie("session")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	userInfo, err := s.db.GetUserFromDB(cookie)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	posts, err := s.db.GetMyPosts(userInfo)
	if err != nil {
		log.Println("My posts handler:", err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	isAuthStruct.AllPosts = posts

	if err := myPostsPage.Execute(w, isAuthStruct); err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
}

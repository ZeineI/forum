package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/ZeineI/forum/internal/models"
)

func (s *Server) FavouritePosts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/favourite" {
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

	favPage, err := template.ParseFiles(files...)
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

	idOfPosts, err := s.db.GetAllIDPostsFav(userInfo)
	if err != nil {
		log.Println("fav handler:", err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	log.Println(idOfPosts)

	posts, err := s.db.GetAllPostsFav(idOfPosts)
	if err != nil {
		log.Println("fav handler:", err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	isAuthStruct.AllPosts = posts

	if err := favPage.Execute(w, isAuthStruct); err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
}

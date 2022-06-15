package server

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"unicode"

	"github.com/ZeineI/forum/internal/models"
)

func (s *Server) Filter(w http.ResponseWriter, r *http.Request) {

	if !strings.HasPrefix(r.URL.Path, "/filter") {
		ErrorHandler(w, http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	files := []string{
		"./frontend/html/mainPage.html",
		"./frontend/html/header.html",
	}

	filterPage, err := template.ParseFiles(files...)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	isAuthStruct := new(models.IsAuthStruct)
	isAuthStruct.IsAuth = s.CheckCookie(w, r)

	tagName := r.FormValue("tag")

	tagAfterTrim := strings.TrimFunc(tagName, func(r rune) bool {
		return unicode.IsSpace(r)
	})

	if tagAfterTrim == "" {
		log.Println("empty tag - miss")
		ErrorHandler(w, http.StatusBadRequest)
		return
	}

	existTag := s.db.IsTagExist(tagName)
	if !existTag {
		if err := filterPage.Execute(w, isAuthStruct); err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		return
	}

	allPosts, err := s.db.GetAllPostsByTag(tagName)
	if err != nil {
		log.Println(err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	isAuthStruct.AllPosts = allPosts

	if err := filterPage.Execute(w, isAuthStruct); err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
}

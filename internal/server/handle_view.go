package server

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/ZeineI/forum/internal/models"
)

func (s *Server) ViewHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/view/") {
		ErrorHandler(w, http.StatusBadRequest)
		return
	}
	if r.Method != http.MethodGet {
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	files := []string{
		"./frontend/html/view.html",
		"./frontend/html/header.html",
	}
	viewPage, err := template.ParseFiles(files...)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	//check cookie
	isAuthStruct := new(models.IsAuthStruct)
	isAuthStruct.IsAuth = s.CheckCookie(w, r)
	id := r.URL.Path[6:]
	idN, err := strconv.Atoi(id)
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest)
		return
	}

	//exist
	existPost := s.db.IsPostExist(idN)
	if !existPost {
		ErrorHandler(w, http.StatusNotFound)
		return
	}
	postInfo, err := s.db.GetPost(idN)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	tags, err := s.db.SelectTag(postInfo.IdPost)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	postInfo.Category = tags
	isAuthStruct.Post = postInfo

	allComments, err := s.db.GetAllComments(postInfo.IdPost)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	isAuthStruct.AllCommentsOnePost = allComments

	//rating
	likes, dislikes, err := s.db.GetAllLikes(postInfo.IdPost)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	isAuthStruct.LikePost = likes
	isAuthStruct.DislikePost = dislikes

	if err := viewPage.Execute(w, isAuthStruct); err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
}

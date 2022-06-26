package server

import (
	"database/sql"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ZeineI/forum/internal/models"
)

func (s *Server) CommentLikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/commentLike" {
		ErrorHandler(w, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	values, err := url.ParseQuery(string(bytes))
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest)
		return
	}
	comment := &models.Comment{}
	for i, v := range values {
		switch i {
		case "CommentLike":
			comment.IdComment, err = strconv.Atoi(v[0])
			if err != nil {
				ErrorHandler(w, http.StatusBadRequest)
				return
			}
		default:
			ErrorHandler(w, http.StatusBadRequest)
			return
		}
	}

	oldUrl := r.Header.Get("Referer")
	idPost, err := strconv.Atoi(oldUrl[27:])
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

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

	newLike := &models.CommentRating{
		Symbol: "like",
		PostId: idPost,
	}
	symbol, err := s.db.CommentDisLikeChecker(newLike, userInfo.Id, comment.IdComment)
	if symbol == "" && err == sql.ErrNoRows {
		if err := s.db.CommentDisLikeInsert(newLike, userInfo, comment); err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
	} else if symbol == "dislike" {
		// if there is a dislike -> update
		if err := s.db.CommentDisLikeUpdate(newLike, userInfo, "like", comment.IdComment); err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
	} else if symbol == "like" {
		// if there is a like already -> delete
		if err := s.db.CommentDisLikeDelete(newLike, userInfo, comment.IdComment); err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, oldUrl, http.StatusSeeOther)
}

func (s *Server) CommentDislikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/commentDislike" {
		ErrorHandler(w, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	values, err := url.ParseQuery(string(bytes))
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest)
		return
	}
	comment := &models.Comment{}
	for i, v := range values {
		switch i {
		case "CommentDislike":
			comment.IdComment, err = strconv.Atoi(v[0])
			if err != nil {
				ErrorHandler(w, http.StatusBadRequest)
				return
			}
		default:
			ErrorHandler(w, http.StatusBadRequest)
			return
		}
	}

	oldUrl := r.Header.Get("Referer")
	idPost, err := strconv.Atoi(oldUrl[27:])
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

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

	newLike := &models.CommentRating{
		Symbol: "dislike",
		PostId: idPost,
	}

	symbol, err := s.db.CommentDisLikeChecker(newLike, userInfo.Id, comment.IdComment)
	if err == sql.ErrNoRows {
		if err := s.db.CommentDisLikeInsert(newLike, userInfo, comment); err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

	} else if symbol == "like" {
		if err := s.db.CommentDisLikeUpdate(newLike, userInfo, "dislike", comment.IdComment); err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

	} else if symbol == "dislike" {
		if err := s.db.CommentDisLikeDelete(newLike, userInfo, comment.IdComment); err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, oldUrl, http.StatusSeeOther)
}

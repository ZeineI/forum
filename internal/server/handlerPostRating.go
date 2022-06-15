package server

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/ZeineI/forum/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

func (s *Server) PostLikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/like" {
		ErrorHandler(w, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	oldUrl := r.Header.Get("Referer")
	idPost, err := strconv.Atoi(oldUrl[27:])
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		log.Println(err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	userInfo, err := s.db.GetUserFromDB(cookie)
	if err != nil {
		log.Println(err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	newLike := &models.PostRating{
		Symbol: "like",
		PostId: idPost,
	}
	symbol, err := s.db.PostDisLikeChecker(newLike, userInfo.Id)
	if symbol == "" && err == sql.ErrNoRows {
		if err := s.db.PostDisLikeInsert(newLike, userInfo); err != nil {
			log.Println(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
	} else if symbol == "dislike" {
		// if there is a dislike -> update
		if err := s.db.PostDisLikeUpdate(newLike, userInfo, "like"); err != nil {
			log.Println(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
	} else if symbol == "like" {
		// if there is a like already -> delete
		if err := s.db.PostDisLikeDelete(newLike, userInfo); err != nil {
			log.Println(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, oldUrl, http.StatusSeeOther)
}

func (s *Server) PostDislikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/dislike" {
		ErrorHandler(w, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	oldUrl := r.Header.Get("Referer")
	log.Println("OLD URL", oldUrl)
	idPost, err := strconv.Atoi(oldUrl[27:])
	log.Println("id post", idPost)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		log.Println(err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	userInfo, err := s.db.GetUserFromDB(cookie)
	if err != nil {
		log.Println(err)
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	log.Println("USER", userInfo)
	newLike := &models.PostRating{
		Symbol: "dislike",
		PostId: idPost,
	}

	symbol, err := s.db.PostDisLikeChecker(newLike, userInfo.Id)
	if err == sql.ErrNoRows {
		if err := s.db.PostDisLikeInsert(newLike, userInfo); err != nil {
			log.Println(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

	} else if symbol == "like" {
		if err := s.db.PostDisLikeUpdate(newLike, userInfo, "dislike"); err != nil {
			log.Println(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

	} else if symbol == "dislike" {
		if err := s.db.PostDisLikeDelete(newLike, userInfo); err != nil {
			log.Println(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, oldUrl, http.StatusSeeOther)
}

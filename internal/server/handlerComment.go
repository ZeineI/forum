package server

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"unicode"

	"github.com/ZeineI/forum/internal/models"
)

func (s *Server) CommentInvisibleHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/postcomment" {
		ErrorHandler(w, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	oldUrl := r.Header.Get("Referer")
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

	idPost, err := strconv.Atoi(oldUrl[27:])
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest)
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
	newComment := &models.Comment{
		PostID: idPost,
	}
	for i, v := range values {
		switch i {
		case "like":
		case "dislike":
		case "submit":
		case "input":
			newComment.TextComment = v[0]
		default:
			ErrorHandler(w, http.StatusBadRequest)
			return
		}
	}

	// check that comment is not empty (not only "required" in HTML), but that not just enters and spaces 	// '\t', '\n', '\v', '\f', '\r', ' ', U+0085 (NEL), U+00A0 (NBSP).
	x := strings.TrimFunc(newComment.TextComment, func(r rune) bool {
		return unicode.IsSpace(r)
	})
	if x == "" {
		log.Println("empty comment - miss")
		ErrorHandler(w, http.StatusBadRequest)
	} else {
		if err := s.db.InsertComment(newComment, userInfo); err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, oldUrl, http.StatusSeeOther)
	}
}

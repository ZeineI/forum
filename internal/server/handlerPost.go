package server

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"unicode"

	"github.com/ZeineI/forum/internal/models"
)

const (
	MB = 1 << 20
)

func (s *Server) PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post" {
		ErrorHandler(w, http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		files := []string{
			"./frontend/html/post.html",
			"./frontend/html/header.html",
		}
		postPage, err := template.ParseFiles(files...)
		if err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		tags := []string{"it", "animal", "random", "anime", "music"}
		auth := &models.IsAuthStruct{
			IsAuth: true,
			Tags:   tags,
		}
		if err := postPage.Execute(w, auth); err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		postInfo := new(models.Post)

		// if you want to limit full request size use:
		r.Body = http.MaxBytesReader(w, r.Body, 32*MB)
		// if largr - it will stop read the request body

		// if you want to limit the amount of memory used use:
		r.ParseMultipartForm(32 * MB)
		// his will use up to maxMemory bytes for file parts, with the remainder stored in temporary files on disk. This call does not limit the total number of bytes read from the client or the size of an uploaded file.

		// https://golangbyexample.com/net-http-package-get-query-params-golang/
		var textPostAfterTrim string
		for i, v := range r.Form {
			switch i {
			case "tags":
				postInfo.Category = v
			case "input":
				postInfo.TextPost = v[0]
				// check that not empty
				textPostAfterTrim = strings.TrimFunc(postInfo.TextPost, func(r rune) bool {
					return unicode.IsSpace(r)
				})
			case "image-post":
			case "submit":
			default:
				ErrorHandler(w, http.StatusBadRequest)
				return
			}
		}

		//upload image
		imageFileName, err := s.uploadFile(w, r)
		if err != nil {
			if err.Error() == "image is larger than 20MB" || err.Error() == "The provided file format is not allowed" {
				log.Println(err)
				ErrorHandler(w, http.StatusBadRequest)
				return
			}
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		postInfo.ImageName = imageFileName

		// tag unique
		for i := 0; i < len(postInfo.Category); i++ {
			for j := 0; j < len(postInfo.Category); j++ {
				if postInfo.Category[i] == postInfo.Category[j] && i != j {
					log.Println("not unique tag")
					ErrorHandler(w, http.StatusBadRequest)
					return
				}
			}
		}
		//tag created by system
		tags := []string{"it", "random", "animal", "anime", "music"}
		for i := 0; i < len(postInfo.Category); i++ {
			counter := 0
			for j := 0; j < len(tags); j++ {
				if postInfo.Category[i] != tags[j] {
					counter++
				}
			}
			if counter == 5 {
				log.Println("no such tag")
				ErrorHandler(w, http.StatusBadRequest)
				return
			}
		}

		// check that at least one is not empty: text post or image
		if textPostAfterTrim == "" && imageFileName == "" {
			log.Println("empty post - miss")
			ErrorHandler(w, http.StatusBadRequest)
		} else {
			if err := s.PostCreation(r, postInfo); err != nil {
				log.Println("Post handler:", err)
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/main", http.StatusSeeOther)
		}
	default:
		ErrorHandler(w, http.StatusMethodNotAllowed)
	}
}

package server

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"text/template"

	"github.com/ZeineI/forum/internal/models"
	"github.com/ZeineI/forum/pkg"
	_ "github.com/mattn/go-sqlite3"
)

// Register
func (s *Server) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// check URL path
	if r.URL.Path != "/register" {
		ErrorHandler(w, http.StatusBadRequest)
		return
	}

	//auth
	IsAuthStruct, err := s.AuthPostOperations(w, r)
	if IsAuthStruct.IsAuth {
		http.Redirect(w, r, "/main", http.StatusSeeOther)
		return
	}

	registerPage, err := template.ParseFiles("./frontend/html/register.html")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	switch r.Method {
	case http.MethodGet:
		if err := registerPage.Execute(w, nil); err != nil {
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		// https://stackoverflow.com/questions/15407719/in-gos-http-package-how-do-i-get-the-query-string-on-a-post-request
		// https://hackthedeveloper.com/golang-forms-data-request-body/- alternative: r.ParseForm and r.PostForm.Get("key")
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
		user := &models.User{}
		for i, v := range values {
			switch i {
			case "username":
				user.Username = v[0]
			case "e-mail":
				user.Email = v[0]
			case "password":
				hash, err := pkg.GeneratePassword(v[0])
				if err != nil {
					log.Println(err)
					ErrorHandler(w, http.StatusInternalServerError)
					return
				}
				user.Password = hash
			case "confirm":
				user.Confirm = v[0]
			default:
				ErrorHandler(w, http.StatusBadRequest)
				return
			}
		}
		user.Place = "server"
		isCorrect, err := CheckEmail(user.Email)
		if err != nil {
			log.Println(err)
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		if !isCorrect {
			w.WriteHeader(400)
			InCorrectRegister(registerPage, w, "email is incorret")
			return
		}
		if !pkg.CheckPasswords(user.Confirm, user.Password) || !IsValidPassword(user.Confirm) {
			w.WriteHeader(400)
			InCorrectRegister(registerPage, w, "password is incorrect")
			log.Println("password is not confirmed or incorrect format")
			return
		}
		if err := s.db.InsertUser(user); err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				clientError := WitchNotUnique(err.Error())
				w.WriteHeader(400)
				InCorrectRegister(registerPage, w, clientError)
				return
			}
			ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		log.Println("registered and redirected")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	default:
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
}

func InCorrectRegister(templ *template.Template, w http.ResponseWriter, clientString string) {
	if err := templ.Execute(w, clientString); err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}
}

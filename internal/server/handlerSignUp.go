package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"text/template"

	models "github.com/ZeineI/forum/internal/models"
	"github.com/dgrijalva/jwt-go"
)

func (s *Server) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// check URL path
	if r.URL.Path != "/register" {
		ErrorHandler(w, http.StatusBadRequest)
		return
	}

	registerPage, err := template.ParseFiles("static/register.html")
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
				user.Password = v[0]
			case "confirm":
				user.Confirm = v[0]
			default:
				ErrorHandler(w, http.StatusBadRequest)
				return
			}
		}
		user.Place = "server"
		if err := s.db.InsertUser(user); err != nil {
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

func IsAuthorized(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] == nil {
			json.NewEncoder(w).Encode(nil)
			http.Redirect(w, r, "/exist", http.StatusSeeOther)
		}
		secretkey := "abc"
		var mySigningKey = []byte(secretkey)

		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error in parsing")
			}
			return mySigningKey, nil
		})

		if err != nil {
			json.NewEncoder(w).Encode(err)
			http.Redirect(w, r, "/exist", http.StatusSeeOther)
		}

		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			handler.ServeHTTP(w, r)
		}
		json.NewEncoder(w).Encode(err)
		http.Redirect(w, r, "/exist", http.StatusSeeOther)
	}
}

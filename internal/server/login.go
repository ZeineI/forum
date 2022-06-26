package server

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/ZeineI/forum/internal/models"
	"github.com/ZeineI/forum/pkg"
	uuid "github.com/satori/go.uuid"
)

// Login
func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// check URL path
	if r.URL.Path != "/login" {
		ErrorHandler(w, http.StatusBadRequest)
		return
	}
	//check auth
	IsAuthStruct, err := s.AuthPostOperations(w, r)
	if IsAuthStruct.IsAuth {
		http.Redirect(w, r, "/main", http.StatusSeeOther)
		return
	}
	loginPage, err := template.ParseFiles("./frontend/html/login.html")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodGet:
		if err := loginPage.Execute(w, nil); err != nil {
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
		userInfo := new(models.User)
		for i, v := range values {
			switch i {
			case "e-mail":
				userInfo.Email = v[0]
			case "password":
				userInfo.Password = v[0]
			default:
				ErrorHandler(w, http.StatusBadRequest)
				return
			}
		}
		userFromDb, err := s.db.GetUser(userInfo.Email)
		if err != nil {
			log.Println("Login: user doesnt exist:", err)
			ErrorHandler(w, http.StatusInternalServerError) //какая ошибка?
			return
		}
		if !pkg.CheckPasswords(userInfo.Password, userFromDb.Password) {
			log.Println("Incorrect password. Refresh page and try again")
			if err := loginPage.Execute(w, "password is inccorect"); err != nil {
				ErrorHandler(w, http.StatusInternalServerError)
				return
			}
			return
		}
		if err := s.Set(w, r, userFromDb); err != nil {
			log.Println("Login set error:", err)
			ErrorHandler(w, http.StatusBadRequest) ///какая ошибка?????
			return
		}
		http.Redirect(w, r, "/main", http.StatusSeeOther)
	default:
		ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
}

func (s *Server) Set(w http.ResponseWriter, r *http.Request, userInfo *models.User) error {
	oldCookie, err := r.Cookie("session")
	if err == nil {
		oldCookie.MaxAge = -1
	}
	http.SetCookie(w, oldCookie)
	id := uuid.NewV4()
	cookie := &http.Cookie{
		Name:     "session",
		Value:    id.String(),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	if err := s.cookies.CheckDoubleCookie(userInfo.Email); err == nil {
		log.Println("Cookie exist")
		userA, err := s.db.GetUser(userInfo.Email)
		if err != nil {
			return fmt.Errorf("Cant get user for update user")
		}
		if errUpdate := s.cookies.UpdateCoockieToBD(cookie.Value, userA); errUpdate != nil {
			log.Println("Login Set Cookie mistake")
			return fmt.Errorf("Update cookie - %w", errUpdate)
		}
		return err
	}
	log.Println("Cookie doesnt exist need to create new")
	if err := s.cookies.AddCoockieToBD(cookie.Value, userInfo); err != nil {
		return fmt.Errorf("Add cookie - %w", err)
	}
	return nil
}

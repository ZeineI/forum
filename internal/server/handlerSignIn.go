package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	models "github.com/ZeineI/forum/internal/models"
	"github.com/dgrijalva/jwt-go"
)

func GenerateJWT(email string) (string, error) {
	secretkey := "abc"
	var mySigningKey = []byte(secretkey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func (s *Server) SignIn(w http.ResponseWriter, r *http.Request) {

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
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	validToken, err := GenerateJWT(userFromDb.Email)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	var token Token
	token.Email = userFromDb.Email
	token.TokenString = validToken
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
	http.Redirect(w, r, "/main", http.StatusSeeOther)
}

type Token struct {
	Role        string `json:"role"`
	Email       string `json:"email"`
	TokenString string `json:"token"`
}

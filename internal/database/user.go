package database

import (
	"fmt"
	"log"
	"net/http"

	models "github.com/ZeineI/forum/internal/models"
)

func (s SqlLiteDB) InsertUser(user *models.User) error {
	result, err := s.DB.Exec("INSERT INTO User(email, username, password, place) VALUES($1, $2, $3, $4)", user.Email, user.Username, user.Password, user.Place)
	if err != nil {
		return err
	}
	log.Println(result.LastInsertId()) // id последнего добавленного объекта
	log.Println(result.RowsAffected()) // количество добавленных строк
	return nil
}

func (s SqlLiteDB) GetUser(email string) (*models.User, error) {
	var (
		usernameDB string
		passWord   string
		id         int
	)
	rows, err := s.DB.Query("SELECT id, username, password FROM User WHERE email=$1", email)
	if err != nil {
		return nil, fmt.Errorf("DB Get User Error (query) - %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&id, &usernameDB, &passWord); err != nil {
			return nil, fmt.Errorf("DB Get User Error (scan) - %w", err)
		}
	}
	user := &models.User{
		Id:       id,
		Email:    email,
		Username: usernameDB,
		Password: passWord,
	}
	return user, nil
}

func (s SqlLiteDB) AlreadyExist(username string) error {
	var usernameGET string
	row := s.DB.QueryRow("SELECT username FROM User WHERE username = $1", username)
	err := row.Scan(&usernameGET)
	log.Println("dbCheckUser:", err)
	return err
}

func (s SqlLiteDB) GetUserForAuth(login string) (*models.User, error) {
	var (
		id int
	)
	rows, err := s.DB.Query("SELECT id FROM User WHERE username=$1", login)
	if err != nil {
		return nil, fmt.Errorf("DB Get User Error (query) - %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("DB Get User Error (scan) - %w", err)
		}
	}
	user := &models.User{
		Id: id,
	}
	return user, nil
}

func (s SqlLiteDB) GetUserFromDB(cookie *http.Cookie) (*models.User, error) {
	user := &models.User{}
	row := s.DB.QueryRow("SELECT User.id, User.username FROM User JOIN Cookies on User.id = user_id WHERE session_id = ? ", cookie.Value)
	if err := row.Scan(&user.Id, &user.Username); err != nil {
		log.Println(err)
		return user, err
	}

	// if err := row.Err(); err != nil {
	// 	return user, err
	// }

	return user, nil
}

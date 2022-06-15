package session

import (
	"fmt"
	"log"

	"github.com/ZeineI/forum/internal/models"
)

func (db CookieStorage) CheckDoubleCookie(email string) error {
	var session_id string
	row := db.db.QueryRow(`SELECT session_id 
	FROM Cookies
	INNER JOIN User
	ON user_id = User.id
	WHERE User.email = $1`, email)
	err := row.Scan(&session_id)
	log.Println(session_id, err)
	return err
}

func (db CookieStorage) CheckDoubleCookieForAuth(login string) error {
	var session_id string
	row := db.db.QueryRow(`SELECT session_id 
	FROM Cookies
	INNER JOIN User
	ON user_id = User.id
	WHERE User.username = $1`, login)
	err := row.Scan(&session_id)
	log.Println(session_id, err)
	return err
}

func (db CookieStorage) AddCoockieToBD(cookieUUID string, userInfo *models.User) error {
	result, err := db.db.Exec("INSERT INTO Cookies(session_id, user_id) VALUES($1, $2)", cookieUUID, userInfo.Id)
	if err != nil {
		return fmt.Errorf("DB Insert Cookie Error - %w", err)
	}
	log.Println(result.LastInsertId())
	log.Println(result.RowsAffected())
	return nil
}

func (s CookieStorage) UpdateCoockieToBD(cookieUUID string, userInfo *models.User) error {
	result, err := s.db.Exec("UPDATE Cookies SET session_id = $1 WHERE user_id = $2", cookieUUID, userInfo.Id)
	if err != nil {
		log.Println("DB Cookie Update Error")
		return fmt.Errorf("DB Cookie Update Error - %w", err)
	}
	log.Println(result.LastInsertId())
	log.Println(result.RowsAffected())
	return nil
}

func (db CookieStorage) DeleteCookie(UUID string) error {
	query := "DELETE FROM Cookies WHERE session_id = $1"
	result, err := db.db.Exec(query, UUID)
	if err != nil {
		return fmt.Errorf("DB Delete Cookie Error - %w", err)
	}
	log.Println(result.LastInsertId())
	log.Println(result.RowsAffected())
	return nil
}

func (db CookieStorage) CheckSession(UUID string) error {
	var session_id string
	row := db.db.QueryRow("SELECT session_id FROM Cookies WHERE session_id = $1", UUID)
	err := row.Scan(&session_id)
	log.Println("dbCheckSession:", err)
	return err
}

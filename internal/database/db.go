package database

import (
	"database/sql"
	"fmt"
	"log"

	models "github.com/ZeineI/forum/internal/models"
)

type SqlLiteDB struct {
	db *sql.DB
}

func (db *SqlLiteDB) Init(dbFile string) (err error) {
	db.db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		return fmt.Errorf("DB: sql open %w", err)
	}

	err = db.db.Ping()
	if err != nil {
		return fmt.Errorf("DB: sql ping %w", err)
	}

	if _, err = db.db.Exec(`CREATE TABLE IF NOT EXISTS User (
		id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		email TEXT UNIQUE, 
		username TEXT UNIQUE NOT NULL, 
		password TEXT,
		place TEXT
	)`); err != nil {
		return fmt.Errorf("DB: User table create %w", err)
	}

	return nil
}

func (s *SqlLiteDB) InsertUser(user *models.User) error {
	result, err := s.db.Exec("INSERT INTO User(email, username, password, place) VALUES($1, $2, $3, $4)", user.Email, user.Username, user.Password, user.Place)
	if err != nil {
		return err
	}
	log.Println(result.LastInsertId()) // id последнего добавленного объекта
	log.Println(result.RowsAffected()) // количество добавленных строк
	return nil
}

func (s *SqlLiteDB) GetUser(email string) (*models.User, error) {
	var (
		usernameDB string
		passWord   string
		id         int
	)
	rows, err := s.db.Query("SELECT id, username, password FROM User WHERE email=$1", email)
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

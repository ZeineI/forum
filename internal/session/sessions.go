package session

import (
	"database/sql"
)

type CookieStorage struct {
	db *sql.DB
}

func InitSession(db *sql.DB) *CookieStorage {
	return &CookieStorage{
		db: db,
	}
}

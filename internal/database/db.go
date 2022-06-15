package database

import (
	"database/sql"
	"fmt"
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

package database

import "database/sql"

type SqlLiteDB struct {
	db *sql.DB
}

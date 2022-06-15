package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type SqlLiteDB struct {
	DB *sql.DB
}

func (db *SqlLiteDB) Init(dbFile string) (err error) {
	db.DB, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		return fmt.Errorf("DB: sql open %w", err)
	}

	err = db.DB.Ping()
	if err != nil {
		return fmt.Errorf("DB: sql ping %w", err)
	}

	if _, err := db.DB.Exec(`CREATE TABLE IF NOT EXISTS User (
		id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		email TEXT UNIQUE, 
		username TEXT UNIQUE NOT NULL, 
		password TEXT,
		place TEXT
	)`); err != nil {
		log.Println("create table(user) error: %v", err)
		return err
	}

	if _, err := db.DB.Exec(`CREATE TABLE IF NOT EXISTS Cookies (
		session_id TEXT,
		user_id TEXT UNIQUE,
		FOREIGN KEY(user_id) REFERENCES User(id) ON DELETE CASCADE
	)`); err != nil {
		log.Println("create table(cookie) error: %v", err)
		return err
	}

	if _, err := db.DB.Exec(`CREATE TABLE IF NOT EXISTS Posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		author_id TEXT NOT NULL,  
		bodyPost TEXT,
		image TEXT,
		FOREIGN KEY(author_Id) REFERENCES User(id) ON DELETE CASCADE,
		CONSTRAINT CK_one_is_not_null CHECK (bodyPost IS NOT NULL OR image IS NOT NULL) 
	)`); err != nil {
		log.Println("create table(post) error: %v", err)
		return err
	}
	// above chnaged - bodypost from Not Null and added imageName
	// added constraint - at least text or image must be there

	if _, err := db.DB.Exec(`CREATE TABLE IF NOT EXISTS Comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, 
		commentAuthor_id INTEGER,
		commentPost_id INTEGER, 
		bodyComment TEXT NOT NULL,
		FOREIGN KEY(commentAuthor_id) REFERENCES User(id) ON DELETE CASCADE,
		FOREIGN KEY(commentPost_id) REFERENCES Posts(id) ON DELETE CASCADE
	)`); err != nil {
		log.Println("create table(comments) error: %v", err)
		return err
	}

	if _, err := db.DB.Exec(`CREATE TABLE IF NOT EXISTS PostRaiting (
		id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		symbol TEXT NOT NULL,
		post_Id INTEGER,
		LikeUser_Id INTEGER, 
		FOREIGN KEY(post_Id) REFERENCES Posts(id) ON DELETE CASCADE,
		FOREIGN KEY(LikeUser_Id) REFERENCES User(id) ON DELETE CASCADE
	)`); err != nil {
		log.Println("create table(postrating) error: %v", err)
		return err
	}

	if _, err := db.DB.Exec(`CREATE TABLE IF NOT EXISTS CommentRaiting (
		id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		symbolComment TEXT NOT NULL,
		commentUser_Id INTEGER,
		comment_Id INTEGER,
		postID INTEGER,
		FOREIGN KEY(commentUser_Id) REFERENCES User(id) ON DELETE CASCADE,
		FOREIGN KEY(comment_Id) REFERENCES Comments(id) ON DELETE CASCADE,
		FOREIGN KEY(postID) REFERENCES Posts(id) ON DELETE CASCADE
	)`); err != nil {
		log.Println("create table(commentrating) error: %v", err)
		return err
	}

	if _, err := db.DB.Exec(`CREATE TABLE IF NOT EXISTS Tags (
		id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		tag TEXT NOT NULL,
		postID INTEGER,
		FOREIGN KEY(postID) REFERENCES Posts(id) ON DELETE CASCADE
	)`); err != nil {
		log.Println("create table(commentrating) error: %v", err)
		return err
	}

	return nil
}

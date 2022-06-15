package database

import (
	"fmt"
	"log"

	models "github.com/ZeineI/forum/internal/models"
)

func (s SqlLiteDB) InitCategories(post *models.Post) error {
	for _, v := range post.Category {
		if err := s.InsertTag(v, post.IdPost); err != nil {
			return err
		}
	}
	return nil
}

func (s SqlLiteDB) InsertTag(tag string, postId int) error {
	query := "INSERT INTO Tags(tag, postID) VALUES($1, $2)"
	result, err := s.DB.Exec(query, tag, postId)
	if err != nil {
		return fmt.Errorf("DB Insert Tag Error - %w", err)
	}
	log.Println(result.LastInsertId())
	log.Println(result.RowsAffected())
	return nil
}

func (s SqlLiteDB) SelectTag(postId int) ([]string, error) {
	var tag string
	var tags []string
	rows, err := s.DB.Query("SELECT tag FROM Tags WHERE postID = $1", postId)
	if err != nil {
		return tags, fmt.Errorf("DB Get Tag Rows Error (query) - %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&tag)
		if err != nil {
			return nil, fmt.Errorf("DB Get Tag One Error (scan) - %w", err)
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (s SqlLiteDB) SelectAllTags() ([]string, error) {
	var tag string
	var tags []string
	rows, err := s.DB.Query("SELECT tag FROM Tags")
	if err != nil {
		return tags, fmt.Errorf("DB Get Tag Rows Error (query) - %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&tag)
		if err != nil {
			return nil, fmt.Errorf("DB Get Tag One Error (scan) - %w", err)
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

package database

import (
	"database/sql"
	"fmt"

	models "github.com/ZeineI/forum/internal/models"
)

func (s SqlLiteDB) IsTagExist(tag string) bool {
	var tagN string
	err := s.DB.QueryRow("SELECT tag FROM Tags WHERE tag = $1", tag).Scan(&tagN)
	if err == sql.ErrNoRows {
		return false
	}
	return true
}

func (s SqlLiteDB) GetAllPostsByTag(tag string) ([]*models.Post, error) {
	var (
		iD        int
		userName  string
		teXt      string
		imageName string
	)
	allPosts := []*models.Post{}
	rows, err := s.DB.Query(`
		SELECT Posts.id, User.username, bodyPost, image
		FROM Posts
		INNER JOIN User ON author_id = User.id
		INNER JOIN Tags ON Posts.id = Tags.postID
		WHERE Tags.tag = $1 
		ORDER BY Posts.id DESC
	`, tag)
	if err != nil {
		return nil, fmt.Errorf("DB Get All Posts Error (query) - %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&iD, &userName, &teXt, &imageName)
		if err != nil {
			return nil, fmt.Errorf("DB Get All Posts Error (scan) - %w", err)
		}
		tags, err := s.SelectTag(iD)
		if err != nil {
			return nil, fmt.Errorf("DB Get All Tags Error (scan) - %w", err)
		}
		allPosts = append(allPosts, &models.Post{
			IdPost:    iD,
			Username:  userName,
			Category:  tags,
			TextPost:  teXt,
			ImageName: imageName,
			Img:       imageName != "",
		})
	}
	return allPosts, nil
}

package database

import (
	"fmt"
	"log"

	models "github.com/ZeineI/forum/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

func (s SqlLiteDB) InsertComment(comment *models.Comment, userInfo *models.User) error {
	result, err := s.DB.Exec("INSERT INTO Comments(commentAuthor_id, commentPost_id, bodyComment) VALUES($1, $2, $3)", userInfo.Id, comment.PostID, comment.TextComment)
	if err != nil {
		return fmt.Errorf("DB Insert Post Error - %w", err)
	}
	log.Println(result.LastInsertId())
	log.Println(result.RowsAffected())
	return nil
}

func (s SqlLiteDB) GetAllComments(postID int) ([]*models.Comment, error) {
	var AllComments []*models.Comment
	var (
		id              int
		usernameComment string
		comment         string
	)
	rows, err := s.DB.Query(`
		SELECT Comments.id, User.username, bodyComment
		FROM Comments
		INNER JOIN User
		ON commentAuthor_id = User.id
		WHERE commentPost_id = $1
	`, postID)
	if err != nil {
		return nil, fmt.Errorf("DB Get All Comments Error (query) - %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &usernameComment, &comment)
		if err != nil {
			return nil, fmt.Errorf("DB Get All Comments Error (scan) - %w", err)
		}
		like, dislike, errComment := s.GetAllLikesComment(postID, id)
		if errComment != nil {
			return AllComments, errComment
		}
		AllComments = append(AllComments, &models.Comment{
			IdComment:   id,
			PostID:      postID,
			Username:    usernameComment,
			TextComment: comment,
			Like:        like,
			Dislike:     dislike,
		})
	}
	return AllComments, nil
}

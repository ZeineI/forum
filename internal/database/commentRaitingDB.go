package database

import (
	"fmt"
	"log"

	models "github.com/ZeineI/forum/internal/models"
)

func (s SqlLiteDB) CommentDisLikeChecker(newLike *models.CommentRating, liker, commentID int) (string, error) {
	var symbol string
	err := s.DB.QueryRow("SELECT symbolComment FROM CommentRaiting WHERE postID = $1 AND commentUser_Id = $2 AND comment_Id = $3", newLike.PostId, liker, commentID).Scan(&symbol)
	return symbol, err
}

func (s SqlLiteDB) CommentDisLikeInsert(newLike *models.CommentRating, user *models.User, comment *models.Comment) error {
	query := "INSERT INTO CommentRaiting(symbolComment, commentUser_Id, comment_Id, postID) VALUES($1, $2, $3, $4)"
	result, err := s.DB.Exec(query, newLike.Symbol, user.Id, comment.IdComment, newLike.PostId)
	if err != nil {
		return fmt.Errorf("DB Insert Comment Like/Dislike Error - %w", err)
	}
	log.Println(result.LastInsertId())
	log.Println(result.RowsAffected())
	return nil
}

func (s SqlLiteDB) CommentDisLikeUpdate(newLike *models.CommentRating, user *models.User, symbolLikeorDislike string, commentId int) error {
	query := "UPDATE CommentRaiting SET symbolComment = $1 WHERE commentUser_Id = $2 AND comment_Id = $3 AND postID = $4"
	result, err := s.DB.Exec(query, symbolLikeorDislike, user.Id, commentId, newLike.PostId)
	if err != nil {
		log.Println("DB Comment Rating Update Error")
		return fmt.Errorf("DB Comment Rating Update Error - %w", err)
	}
	log.Println(result.LastInsertId())
	log.Println(result.RowsAffected())
	return nil
}

func (s SqlLiteDB) CommentDisLikeDelete(newLike *models.CommentRating, user *models.User, commentId int) error {
	query := "DELETE FROM CommentRaiting WHERE postID = $1 AND commentUser_Id = $2 AND comment_Id = $3"
	result, err := s.DB.Exec(query, newLike.PostId, user.Id, commentId)
	if err != nil {
		return fmt.Errorf("DB Delete Comment Rating Error - %w", err)
	}
	log.Println(result.LastInsertId())
	log.Println(result.RowsAffected())
	return nil
}

func (s SqlLiteDB) GetAllLikesComment(postId, commentId int) (int, int, error) {
	var (
		symbol  string
		like    int
		dislike int
	)
	rows, err := s.DB.Query("SELECT symbolComment FROM CommentRaiting WHERE postID = $1 AND comment_Id = $2", postId, commentId)
	if err != nil {
		return 0, 0, fmt.Errorf("DB Get All Likes Error (query) - %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&symbol)
		if err != nil {
			return 0, 0, fmt.Errorf("DB Get All Likes Error (scan) - %w", err)
		}
		if symbol == "like" {
			like++
		} else {
			dislike++
		}
	}
	return like, dislike, nil
}

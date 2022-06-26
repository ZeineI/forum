package database

import (
	"fmt"
	"log"

	models "github.com/ZeineI/forum/internal/models"
)

// Checks if there is already like or dislike on this post by this user
func (s SqlLiteDB) PostDisLikeChecker(newLike *models.PostRating, liker int) (string, error) {
	var symbol string
	err := s.DB.QueryRow("SELECT symbol FROM PostRaiting WHERE post_id = $1 AND LikeUser_Id = $2", newLike.PostId, liker).Scan(&symbol)
	return symbol, err
}

//insert like or dislike
func (s SqlLiteDB) PostDisLikeInsert(newLike *models.PostRating, user *models.User) error {
	query := "INSERT INTO PostRaiting(symbol, post_Id, LikeUser_Id) VALUES($1, $2, $3)"
	result, err := s.DB.Exec(query, newLike.Symbol, newLike.PostId, user.Id)
	if err != nil {
		return fmt.Errorf("DB Insert Like/Dislike Error - %w", err)
	}
	log.Println(result.LastInsertId())
	log.Println(result.RowsAffected())
	return nil
}

//update like or dislike
func (s SqlLiteDB) PostDisLikeUpdate(newLike *models.PostRating, user *models.User, symbolLikeorDislike string) error {
	query := "UPDATE Postraiting SET symbol = $1 WHERE LikeUser_Id = $2 AND post_Id = $3"
	result, err := s.DB.Exec(query, symbolLikeorDislike, user.Id, newLike.PostId)
	if err != nil {
		log.Println("DB Post Rating Update Error")
		return fmt.Errorf("DB Post Rating Update Error - %w", err)
	}
	log.Println(result.LastInsertId())
	log.Println(result.RowsAffected())
	return nil
}

//delete like or dislike
func (s SqlLiteDB) PostDisLikeDelete(newLike *models.PostRating, user *models.User) error {
	query := "DELETE FROM PostRaiting WHERE LikeUser_Id = $1 AND post_Id = $2"
	result, err := s.DB.Exec(query, user.Id, newLike.PostId)
	if err != nil {
		return fmt.Errorf("DB Delete Post Rating Error - %w", err)
	}
	log.Println(result.LastInsertId())
	log.Println(result.RowsAffected())
	return nil
}

func (s SqlLiteDB) GetAllLikes(postId int) (int, int, error) {
	var (
		symbol  string
		like    int
		dislike int
	)
	rows, err := s.DB.Query("SELECT symbol FROM PostRaiting WHERE post_Id = $1", postId)
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

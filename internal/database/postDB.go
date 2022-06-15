package database

import (
	"database/sql"
	"fmt"

	models "github.com/ZeineI/forum/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

func (s SqlLiteDB) InsertPost(user *models.User, post *models.Post) error {
	result, err := s.DB.Exec("INSERT INTO Posts(author_id, bodyPost, image) VALUES($1, $2, $3)", user.Id, post.TextPost, post.ImageName)
	if err != nil {
		return fmt.Errorf("DB Insert Post Error - %w", err)
	}
	postId, err := result.LastInsertId() // id последнего добавленного объекта
	if err != nil {
		return err
	}
	post.IdPost = int(postId)
	if err := s.InitCategories(post); err != nil {
		_, err := s.DB.Exec("DELETE FROM Posts WHERE Id = $1", post.IdPost)
		if err != nil {
			return fmt.Errorf("DB Insert Category Error - %w", err)
		}
	}
	return nil
}

func (s SqlLiteDB) GetAllPosts() ([]*models.Post, error) {
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
		INNER JOIN User 
		ON author_id = User.id
		ORDER BY Posts.id DESC
	`)
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

func (s SqlLiteDB) GetPost(id int) (*models.Post, error) {
	var (
		iD        int
		userName  string
		teXt      string
		imageName string
	)
	postInfo := &models.Post{}
	row := s.DB.QueryRow(`
		SELECT Posts.id, User.username, bodyPost, image
		FROM Posts 
		INNER JOIN User 
		ON author_id = User.id 
		WHERE Posts.id = $1
	`, id)
	err := row.Scan(&iD, &userName, &teXt, &imageName)
	if err != nil {
		return postInfo, err
	}
	tags, err := s.SelectTag(iD)
	if err != nil {
		return nil, fmt.Errorf("DB Get All Tags Error (scan) - %w", err)
	}
	postInfo = &models.Post{
		IdPost:    iD,
		Username:  userName,
		Category:  tags,
		TextPost:  teXt,
		ImageName: imageName,
		Img:       imageName != "",
	}
	return postInfo, nil
}

func (s SqlLiteDB) IsPostExist(id int) bool {
	var idN int
	err := s.DB.QueryRow("SELECT id FROM Posts WHERE id = $1", id).Scan(&idN)
	if err == sql.ErrNoRows {
		return false
	}
	return true
}

func (s SqlLiteDB) GetAllIDPostsFav(user *models.User) ([]int, error) {
	var (
		idS []int
		id  int
	)
	rows, err := s.DB.Query(`
		SELECT post_Id
		FROM PostRaiting 
		WHERE LikeUser_Id = $1 AND symbol = $2
		ORDER BY PostRaiting.id DESC
	`, user.Id, "like")
	if err != nil {
		return nil, fmt.Errorf("DB Get All Posts  FAV Error (query) - %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("DB Get All Posts FAV Error (scan) - %w", err)
		}
		idS = append(idS, id)
	}
	return idS, nil
}

func (s SqlLiteDB) GetAllPostsFav(idS []int) ([]*models.Post, error) {
	allPosts := []*models.Post{}
	for _, v := range idS {
		post, err := s.GetPost(v)
		if err != nil {
			return allPosts, fmt.Errorf("dbGetAllPostsFav - %w", err)
		}
		allPosts = append(allPosts, post)
	}
	return allPosts, nil
}

func (s SqlLiteDB) GetMyPosts(user *models.User) ([]*models.Post, error) {
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
		INNER JOIN User 
		ON author_id = User.id
		WHERE User.id = $1
		ORDER BY Posts.id DESC
	`, user.Id)
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

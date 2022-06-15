package database

import (
	"net/http"

	models "github.com/ZeineI/forum/internal/models"
)

type Storage interface {
	User
	Post
	Comments
	Categories
}

type Comments interface {
	// comments
	InsertComment(comment *models.Comment, userInfo *models.User) error
	GetAllComments(postID int) ([]*models.Comment, error)
	// like comments
	CommentDisLikeChecker(newLike *models.CommentRating, liker, commentID int) (string, error)
	CommentDisLikeInsert(newLike *models.CommentRating, user *models.User, comment *models.Comment) error
	CommentDisLikeUpdate(newLike *models.CommentRating, user *models.User, symbolLikeorDislike string, commentId int) error
	CommentDisLikeDelete(newLike *models.CommentRating, user *models.User, commentId int) error
	GetAllLikesComment(postId, commentId int) (int, int, error)
}

type Categories interface {
	// categories
	InitCategories(post *models.Post) error
	InsertTag(tag string, postId int) error
	SelectTag(postId int) ([]string, error)
	SelectAllTags() ([]string, error)

	//filter
	IsTagExist(tag string) bool

	GetAllPostsByTag(tag string) ([]*models.Post, error)
}

type Post interface {
	InsertPost(user *models.User, post *models.Post) error
	GetAllPosts() ([]*models.Post, error)
	GetPost(id int) (*models.Post, error)
	IsPostExist(id int) bool
	GetAllIDPostsFav(user *models.User) ([]int, error)
	GetAllPostsFav(idS []int) ([]*models.Post, error)
	GetMyPosts(user *models.User) ([]*models.Post, error)
	// like posts
	PostDisLikeChecker(newLike *models.PostRating, liker int) (string, error)
	PostDisLikeInsert(newLike *models.PostRating, user *models.User) error
	PostDisLikeUpdate(newLike *models.PostRating, user *models.User, symbolLikeorDislike string) error
	PostDisLikeDelete(newLike *models.PostRating, user *models.User) error
	GetAllLikes(postId int) (int, int, error)
}

type User interface {
	InsertUser(user *models.User) error
	// GetPlace(session string) (string, error)
	AlreadyExist(username string) error
	GetUser(email string) (*models.User, error)
	GetUserForAuth(login string) (*models.User, error)
	GetUserFromDB(cookie *http.Cookie) (*models.User, error)
}

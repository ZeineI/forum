package models

type User struct {
	Id       int
	Email    string
	Username string
	Password string
	Confirm  string
	Place    string
}

type Post struct {
	IdPost    int
	Username  string //connection to user for database (user who created the post)
	Category  []string
	TextPost  string //text of post
	ImageName string
	Img       bool
}

type IsAuthStruct struct {
	AllPosts           []*Post
	Tags               []string
	Post               *Post
	AllCommentsOnePost []*Comment
	IsAuth             bool
	LikePost           int
	DislikePost        int
}

type Comment struct {
	IdComment   int    //own id
	PostID      int    //connection to post (comment on which post)
	Username    string //connection to user for database (user who commented)
	TextComment string //text of comment
	Like        int
	Dislike     int
}

type ErrorStruct struct {
	ErrorNum int
	CodeText string
}

type PostRating struct {
	IdLike   int
	Symbol   string //true - like, false - dislike
	Username string
	PostId   int
}

type CommentRating struct {
	IdLike    int
	Symbol    string //true - like, false - dislike
	Username  string
	PostId    int
	CommentId int
}

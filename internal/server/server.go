package server

import (
	"fmt"
	"net/http"

	"github.com/ZeineI/forum/internal/database"
	"github.com/ZeineI/forum/internal/session"
)

type Server struct {
	mux     *http.ServeMux
	db      database.Storage
	cookies session.CookieRepository
}

func InitServer(db database.Storage, cookies session.CookieRepository) *Server {
	return &Server{
		mux:     http.NewServeMux(),
		db:      db,
		cookies: cookies,
	}
}

func (s *Server) Run() error {
	s.mux.Handle("/frontend/", http.StripPrefix("/frontend/", http.FileServer(http.Dir("./frontend"))))

	//HANDLERS
	s.mux.HandleFunc("/main", s.MainHandler)
	s.mux.HandleFunc("/register", s.RegisterHandler)
	s.mux.HandleFunc("/login", s.LoginHandler)

	s.mux.Handle("/post", s.MiddlewareLoginCheck(http.HandlerFunc(s.PostHandler)))
	s.mux.Handle("/postcomment", s.MiddlewareLoginCheck(http.HandlerFunc(s.CommentInvisibleHandler)))
	s.mux.HandleFunc("/view/", s.ViewHandler)

	//rating
	s.mux.Handle("/like", s.MiddlewareLoginCheck(http.HandlerFunc(s.PostLikeHandler)))
	s.mux.Handle("/dislike", s.MiddlewareLoginCheck(http.HandlerFunc(s.PostDislikeHandler)))

	//comment rating
	s.mux.Handle("/commentDislike", s.MiddlewareLoginCheck(http.HandlerFunc(s.CommentDislikeHandler)))
	s.mux.Handle("/commentLike", s.MiddlewareLoginCheck(http.HandlerFunc(s.CommentLikeHandler)))

	//add
	s.mux.Handle("/favourite", s.MiddlewareLoginCheck(http.HandlerFunc(s.FavouritePosts)))
	s.mux.Handle("/myPosts", s.MiddlewareLoginCheck(http.HandlerFunc(s.MyPosts)))

	//filter
	s.mux.HandleFunc("/filter", s.Filter)

	//logout
	s.mux.HandleFunc("/logout", s.LogOut)

	//if no such page
	s.mux.HandleFunc("/", s.Check)
	fmt.Printf("A server is running on this address: http://localhost:8081/\n")
	if err := http.ListenAndServe(":8081", s.mux); err != nil {
		return err
	}
	return nil
}

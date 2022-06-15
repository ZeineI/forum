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
	// s.mux.HandleFunc("/", MainHandler)
	fmt.Printf("A server is running on this address: http://localhost:8081/\n")
	if err := http.ListenAndServe(":8081", s.mux); err != nil {
		return err
	}
	return nil
}

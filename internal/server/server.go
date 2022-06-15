package server

import (
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

package server

import (
	"log"
	"net/http"

	"github.com/ZeineI/forum/internal/models"
)

func (s *Server) PostCreation(r *http.Request, postInfo *models.Post) error {
	cookie, err := r.Cookie("session")
	if err != nil {
		return err
	}
	user, err := s.db.GetUserFromDB(cookie)
	if err != nil {
		log.Println(err)
		return err
	}
	if err := s.db.InsertPost(user, postInfo); err != nil {
		return err
	}
	return nil
}

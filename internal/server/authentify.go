package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ZeineI/forum/internal/models"
)

func (s *Server) AuthPostOperations(w http.ResponseWriter, r *http.Request) (*models.IsAuthStruct, error) {
	isAuthStruct := new(models.IsAuthStruct)
	isAuthStruct.IsAuth = s.CheckCookie(w, r)
	switch r.URL.Path {
	case "/main":
		allPosts, err := s.db.GetAllPosts()
		if err != nil {
			return nil, fmt.Errorf("AuthPostOperation main case err: %w")
		}
		isAuthStruct.AllPosts = allPosts
		allDoublTags, err := s.db.SelectAllTags()
		if err != nil {
			return nil, fmt.Errorf("AuthPostOperation all tags err: %w")
		}
		isAuthStruct.Tags = allDoublTags
		// isAuthStruct.Tags = s.RemoveDoublicates(allDoublTags)
	case "/login":
		return isAuthStruct, nil
	case "/register":
		return isAuthStruct, nil
	default:
		return nil, fmt.Errorf("AuthPostOperation default path err: %w")
	}

	for _, p := range isAuthStruct.AllPosts {
		log.Println(p.ImageName)
	}

	return isAuthStruct, nil
}

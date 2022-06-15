package session

import (
	"database/sql"
	"net/http"
)

func (s CookieStorage) IsUserInSession(cookie *http.Cookie) bool {
	if cookie == nil {
		return false
	}
	var val string
	err := s.db.QueryRow("SELECT Value FROM Cookie WHERE Value = ?", cookie.Value).Scan(&val)
	if err == nil && err == sql.ErrNoRows {
		return false
	}
	if val == "" {
		return false
	}
	return true
}

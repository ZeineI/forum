package session

import (
	"net/http"

	"github.com/ZeineI/forum/internal/models"
)

type CookieRepository interface {
	IsUserInSession(cookie *http.Cookie) bool
	CheckDoubleCookie(email string) error
	CheckDoubleCookieForAuth(login string) error
	AddCoockieToBD(cookieUUID string, userInfo *models.User) error
	UpdateCoockieToBD(cookieUUID string, userInfo *models.User) error
	DeleteCookie(UUID string) error
	CheckSession(UUID string) error
}

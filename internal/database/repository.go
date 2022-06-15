package database

import (
	models "github.com/ZeineI/forum/internal/models"
)

type Storage interface {
	InsertUser(u *models.User) error
	GetUser(email string) (*models.User, error)
}

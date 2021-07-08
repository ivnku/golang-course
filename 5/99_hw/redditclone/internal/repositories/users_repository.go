package repositories

import (
	"redditclone/internal"
	"redditclone/internal/models"
)

type UsersRepository struct {}

func (ur *UsersRepository) List() ([]*models.User, error) {
	db := internal.DB

	users := []*models.User{}

	_ = db.Find(&users)
	err := db.Error

	if err != nil {
		return users, err
	}

	return users, nil
}

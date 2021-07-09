package user

import (
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) List() ([]*User, error) {
	var users []*User

	db := r.DB.Find(&users)

	if err := db.Error; err != nil {
		return users, err
	}

	return users, nil
}

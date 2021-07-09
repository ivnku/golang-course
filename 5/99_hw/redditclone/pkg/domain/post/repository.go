package post

import (
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) List() ([]*Post, error) {
	var posts []*Post

	db := r.DB.Find(&posts)

	if err := db.Error; err != nil {
		return posts, err
	}

	return posts, nil
}
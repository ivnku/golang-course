package comment

import "gorm.io/gorm"

type Repository struct {
	DB *gorm.DB
}

/**
 * @Description: Create a new post's comment
 * @receiver r
 * @param postComment
 * @return *Comment
 * @return error
 */
func (r *Repository) Create(postComment *Comment) (*Comment, error) {
	db := r.DB.Create(postComment)

	if err := db.Error; err != nil {
		return nil, err
	}

	return postComment, nil
}
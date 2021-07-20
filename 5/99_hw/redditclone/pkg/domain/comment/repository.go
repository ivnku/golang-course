package comment

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{db}
}

/**
 * @Description: Create a new post's comment
 * @receiver r
 * @param postComment
 * @return *Comment
 * @return error
 */
func (r *Repository) Create(postComment *Comment) (*Comment, error) {
	db := r.db.Create(postComment)

	if err := db.Error; err != nil {
		return nil, err
	}

	return postComment, nil
}

/**
 * @Description: Delete a post's comment
 * @receiver r
 * @param id
 * @return bool
 * @return error
 */
func (r *Repository) Delete(id uint) (bool, error) {

	db := r.db.Delete(&Comment{}, id)

	if err := db.Error; err != nil {
		return false, err
	}

	return true, nil
}

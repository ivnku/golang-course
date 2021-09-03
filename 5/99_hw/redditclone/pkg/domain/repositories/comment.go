package repositories

import (
	"gorm.io/gorm"
	"redditclone/pkg/domain/models"
)

type CommentsRepository struct {
	db *gorm.DB
}

func NewCommentsRepository(db *gorm.DB) CommentsRepository {
	return CommentsRepository{db}
}

/**
 * @Description: Create a new post's comment
 * @receiver r
 * @param postComment
 * @return *Comment
 * @return error
 */
func (r *CommentsRepository) Create(postComment *models.Comment) (*models.Comment, error) {
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
func (r *CommentsRepository) Delete(id uint) (bool, error) {

	db := r.db.Delete(&models.Comment{}, id)

	if err := db.Error; err != nil {
		return false, err
	}

	return true, nil
}

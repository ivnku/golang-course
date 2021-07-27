package repositories

import (
	"go.mongodb.org/mongo-driver/mongo"
	"redditclone/pkg/domain/models"
)

type CommentsRepository struct {
	mongodb    *mongo.Client
	collection *mongo.Collection
}

func NewCommentsRepository(db *mongo.Client) CommentsRepository {
	collection := db.Database("redditclone").Collection("comments")
	return CommentsRepository{db, collection}
}

/**
 * @Description: Create a new post's comment
 * @receiver r
 * @param postComment
 * @return *Comment
 * @return error
 */
func (r *CommentsRepository) Create(postComment *models.Comment) (*models.Comment, error) {
	//db := r.db.Create(postComment)
	//
	//if err := db.Error; err != nil {
	//	return nil, err
	//}

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

	//db := r.db.Delete(&models.Comment{}, id)
	//
	//if err := db.Error; err != nil {
	//	return false, err
	//}

	return true, nil
}

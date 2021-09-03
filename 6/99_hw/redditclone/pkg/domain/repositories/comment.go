package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"redditclone/pkg/domain/models"
	"redditclone/pkg/domain/repositories/mocks"
)

type CommentsRepository struct {
	collection mocks.IMongoCollection
}

type ICommentsRepository interface {
	Create(postComment *models.Comment) (*models.Comment, error)
	Delete(id string) (bool, error)
}

func NewCommentsRepository(collection mocks.IMongoCollection) *CommentsRepository {
	return &CommentsRepository{collection}
}

/**
 * @Description: Create a new post's comment
 * @receiver r
 * @param postComment
 * @return *Comment
 * @return error
 */
func (r *CommentsRepository) Create(postComment *models.Comment) (*models.Comment, error) {
	var ctx = context.Background()
	postComment.ID = primitive.NewObjectID()
	result, err := r.collection.InsertOne(ctx, postComment)

	if err != nil {
		return nil, err
	}

	postComment.ID = result.InsertedID.(primitive.ObjectID)

	return postComment, nil
}

/**
 * @Description: Delete a post's comment
 * @receiver r
 * @param id
 * @return bool
 * @return error
 */
func (r *CommentsRepository) Delete(id string) (bool, error) {
	var ctx = context.Background()
	commentId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return false, err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": commentId})

	if err != nil {
		return false, err
	}

	return true, nil
}

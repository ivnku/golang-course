package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"redditclone/pkg/domain/models"
	"redditclone/pkg/domain/repositories/mocks"
)

type VotesRepository struct {
	collection mocks.IMongoCollection
}

type IVotesRepository interface {
	Create(vote *models.Vote) (*models.Vote, error)
	Delete(id primitive.ObjectID) (bool, error)
}

func NewVotesRepository(collection mocks.IMongoCollection) *VotesRepository {
	return &VotesRepository{collection}
}

/**
 * @Description: Create a new vote
 * @receiver r
 * @param vote
 * @return *Vote
 * @return error
 */
func (r *VotesRepository) Create(vote *models.Vote) (*models.Vote, error) {
	var ctx = context.Background()
	vote.ID = primitive.NewObjectID()
	result, err := r.collection.InsertOne(ctx, vote)

	if err != nil {
		return nil, err
	}

	vote.ID = result.InsertedID.(primitive.ObjectID)

	return vote, nil
}

/**
 * @Description: Delete the vote
 * @receiver r *PostsRepository
 * @return bool
 * @return error
 */
func (r *VotesRepository) Delete(id primitive.ObjectID) (bool, error) {
	var ctx = context.Background()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})

	if err != nil {
		return false, err
	}

	return true, nil
}

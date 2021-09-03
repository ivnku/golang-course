package tests

import (
	"context"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"redditclone/pkg/domain/models"
	"redditclone/pkg/domain/repositories"
	"redditclone/pkg/domain/repositories/mocks"
	"reflect"
	"testing"
)

/**
 * @Description: Test creation of a vote
 * @param t
 */
func TestVotesRepoCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockCollection := mocks.NewMockIMongoCollection(ctrl)

	repo := repositories.NewVotesRepository(mockCollection)
	objID := primitive.NewObjectID()

	voteToCreate := &models.Vote{
		ID:     objID,
		PostId: objID,
		UserId: 16,
		Vote:   1,
	}

	mockCollection.EXPECT().
		InsertOne(ctx, voteToCreate).
		Return(&mongo.InsertOneResult{InsertedID: objID}, nil)

	res, err := repo.Create(voteToCreate)
	if !reflect.DeepEqual(res, voteToCreate) {
		t.Errorf("bad result, expected %v, got %v", voteToCreate, res)
	}

	if err != nil {
		t.Errorf("unexpected error, got %v", err)
	}
}

/**
 * @Description: Test deletion of a vote
 * @param t
 */
func TestVotesRepoDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockCollection := mocks.NewMockIMongoCollection(ctrl)
	repo := repositories.NewVotesRepository(mockCollection)

	id := "6106df17322d6f1d06993787"
	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		t.Errorf("unexpected error, got %v", err)
	}

	filter := bson.M{"_id": objID}

	mockCollection.EXPECT().
		DeleteOne(ctx, filter).
		Return(&mongo.DeleteResult{DeletedCount: 1}, nil)

	res, err := repo.Delete(objID)

	if err != nil {
		t.Errorf("unexpected error, got %v", err)
	}

	if !res {
		t.Errorf("couldn't remove object!")
	}
}

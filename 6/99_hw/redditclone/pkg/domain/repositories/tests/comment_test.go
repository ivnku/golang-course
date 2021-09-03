package tests

import (
	"context"
	gomock "github.com/golang/mock/gomock"
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
 * @Description: Test creation of a comment
 * @param t
 */
func TestCommentsRepoCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockCollection := mocks.NewMockIMongoCollection(ctrl)

	repo := repositories.NewCommentsRepository(mockCollection)
	objID := primitive.NewObjectID()

	commentToCreate := &models.Comment{
		ID:     objID,
		PostID: primitive.ObjectID{2},
		UserID: 2,
		Body:   "Some comment",
	}

	mockCollection.EXPECT().
		InsertOne(ctx, commentToCreate).
		Return(&mongo.InsertOneResult{InsertedID: objID}, nil)

	res, err := repo.Create(commentToCreate)
	if !reflect.DeepEqual(res, commentToCreate) {
		t.Errorf("bad result, expected %v, got %v", commentToCreate, res)
	}

	if err != nil {
		t.Errorf("unexpected error, got %v", err)
	}
}

/**
 * @Description: Test deletion of a comment
 * @param t
 */
func TestCommentsRepoDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockCollection := mocks.NewMockIMongoCollection(ctrl)
	repo := repositories.NewCommentsRepository(mockCollection)

	id := "6106df17322d6f1d06993787"
	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		t.Errorf("unexpected error, got %v", err)
	}

	filter := bson.M{"_id": objID}

	mockCollection.EXPECT().
		DeleteOne(ctx, filter).
		Return(&mongo.DeleteResult{DeletedCount: 1}, nil)

	res, err := repo.Delete(id)

	if err != nil {
		t.Errorf("unexpected error, got %v", err)
	}

	if !res {
		t.Errorf("couldn't remove object!")
	}
}

/**
 * @Description: Test deletion of a comment with invalid id
 * @param t
 */
func TestCommentsRepoDeleteWithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockIMongoCollection(ctrl)
	repo := repositories.NewCommentsRepository(mockCollection)

	id := "6"

	res, err := repo.Delete(id)

	if !reflect.DeepEqual(primitive.ErrInvalidHex, err) {
		t.Errorf("unexpected error, got %v", err)
	}

	if res {
		t.Errorf("Delete was ended up successfully but fail was expected!")
	}
}

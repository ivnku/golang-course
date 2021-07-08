package main

import (
	"context"
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"

	gomock "github.com/golang/mock/gomock"
)

func TestRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockCollection := NewMockIMongoCollection(ctrl)
	mockSingleResult := NewMockIMongoSingleResult(ctrl)

	repo := &PostsRepo{
		Collection: mockCollection,
	}

	objID := primitive.NewObjectID()

	// cond := bson.M{"_id": objID}
	// cond := bson.D{{"_id", objID}}

	expectedPost := &Post{
		ID:    objID,
		Email: "vasya@mail.ru",
	}

	/*
		err := repo.Collection.FindOne(ctx, bson.M{"_id": objID})
		.Decode(post)
	*/

	mockCollection.EXPECT().
		FindOne(ctx, gomock.Any()).
		Return(mockSingleResult)
	mockSingleResult.EXPECT().
		Decode(gomock.AssignableToTypeOf(expectedPost)).
		// Decode(gomock.Any()).
		SetArg(0, *expectedPost).
		Return(nil)

	res, err := repo.GetOne(ctx, objID)
	if !reflect.DeepEqual(res, expectedPost) {
		t.Errorf("bad result, expected %v, got %v", expectedPost, res)
	}

	if err != nil {
		t.Errorf("unexpected error, got %v", err)
	}
}

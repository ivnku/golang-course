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
 * @Description: Test creation of a post
 * @param t
 */
func TestPostsRepoCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockCollection := mocks.NewMockIMongoCollection(ctrl)

	repo := repositories.NewPostsRepository(mockCollection)
	objID := primitive.NewObjectID()

	postToCreate := &models.Post{
		ID:               objID,
		Title:            "Post title",
		Category:         "Some category",
		Score:            0,
		Type:             "text",
		Url:              "",
		Text:             "post content",
		UpvotePercentage: 0,
		Views:            0,
	}

	mockCollection.EXPECT().
		InsertOne(ctx, postToCreate).
		Return(&mongo.InsertOneResult{InsertedID: objID}, nil)

	res, err := repo.Create(postToCreate)
	if !reflect.DeepEqual(res, postToCreate) {
		t.Errorf("bad result, expected %v, got %v", postToCreate, res)
	}

	if err != nil {
		t.Errorf("unexpected error, got %v", err)
	}
}

/**
 * @Description: Test getting of a list of posts
 * @param t
 */
func TestPostsRepoList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockCollection := mocks.NewMockIMongoCollection(ctrl)
	mockCursor := mocks.NewMockIMongoCursor(ctrl)

	repo := repositories.NewPostsRepository(mockCollection)

	postId1, err := primitive.ObjectIDFromHex("61059ec7adc529aef3ac11b3")

	if err != nil {
		t.Errorf("unexpected error, got %v", err)
	}

	postId2, err := primitive.ObjectIDFromHex("6106e2c85012868c20f37691")

	if err != nil {
		t.Errorf("unexpected error, got %v", err)
	}

	expectedPosts := []*models.Post{
		{
			ID:               postId1,
			Title:            "The first post",
			Category:         "music",
			Score:            0,
			Type:             "text",
			Url:              "",
			Text:             "",
			UpvotePercentage: 0,
			Views:            0,
		},
		{
			ID:               postId2,
			Title:            "The second post",
			Category:         "news",
			Score:            0,
			Type:             "text",
			Url:              "",
			Text:             "",
			UpvotePercentage: 0,
			Views:            0,
		},
	}

	mockCollection.EXPECT().
		Aggregate(ctx, gomock.Any()).
		Return(mockCursor, nil)
	mockCursor.EXPECT().
		All(ctx, &expectedPosts).
		Return(nil)

	mockCursor.EXPECT().Close(gomock.Any())

	res, err := repo.List()
	if len(res) != len(expectedPosts) {
		t.Errorf("bad result, expected %v, got %v", expectedPosts, res)
	}

	if err != nil {
		t.Errorf("unexpected error, got %v", err)
	}
}

/**
 * @Description: Test update of a post
 * @param t
 */
func TestPostsRepoUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockCollection := mocks.NewMockIMongoCollection(ctrl)
	mockUpdateResult := mocks.NewMockIMongoUpdateResult(ctrl)

	repo := repositories.NewPostsRepository(mockCollection)

	objID := primitive.NewObjectID()
	originalPost := &models.Post{
		ID:               objID,
		Title:            "Post title",
		Category:         "Some category",
		Score:            0,
		Type:             "text",
		Url:              "",
		Text:             "post content",
		UpvotePercentage: 0,
		Views:            0,
	}

	updatedPost := &models.Post{
		ID:               objID,
		Title:            "Post title upd",
		Category:         "Some category",
		Score:            0,
		Type:             "text",
		Url:              "",
		Text:             "post content",
		UpvotePercentage: 0,
		Views:            0,
	}

	mockCollection.EXPECT().
		UpdateOne(ctx, bson.M{"_id": originalPost.ID}, primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "Title", Value: "Post title upd"}}}).
		Return(mockUpdateResult, nil)

	res, err := repo.Update(originalPost, []primitive.E{{"Title", "Post title upd"}})
	if !reflect.DeepEqual(res, updatedPost) {
		t.Errorf("bad result, expected %v, got %v", updatedPost, res)
	}

	if err != nil {
		t.Errorf("unexpected error, got %v", err)
	}
}

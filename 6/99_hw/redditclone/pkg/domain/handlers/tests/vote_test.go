package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http/httptest"
	"redditclone/configs"
	"redditclone/pkg/auth"
	"redditclone/pkg/domain/handlers"
	"redditclone/pkg/domain/models"
	"redditclone/pkg/domain/repositories/mocks"
	"testing"
)

func TestVotesHandlerUpvote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := configs.Config{Token: "sometoken"}

	votesRepository := mocks.NewMockIVotesRepository(ctrl)
	postsRepository := mocks.NewMockIPostsRepository(ctrl)
	votesHandler := &handlers.VotesHandler{
		VotesRepository: votesRepository,
		PostsRepository: postsRepository,
		Config:          config,
	}

	postIdString := "qwertyxyzdgrqwertyxyzdgr"
	postId, err := primitive.ObjectIDFromHex(postIdString)
	voteId := primitive.NewObjectID()

	post := &models.Post{
		ID:               postId,
		Title:            "Some post",
		Category:         "music",
		Score:            0,
		Type:             "text",
		Url:              "",
		Text:             "some content",
		UpvotePercentage: 0,
		Views:            0,
	}
	postAfterUpvote := &models.Post{
		ID:               postId,
		Title:            "Some post",
		Category:         "music",
		Score:            1,
		Type:             "text",
		Url:              "",
		Text:             "some content",
		UpvotePercentage: 100,
		Views:            0,
		Votes:            []*models.Vote{},
	}

	// Test upvote for invalid user id
	userData := auth.UserData{
		Id:       "invalidId",
		Username: "User",
	}

	req := httptest.NewRequest("GET", "/api/post/14/upvote", nil)
	ctx := req.Context()
	ctx = context.WithValue(ctx, configs.UserCtx, userData)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	votesHandler.Upvote(w, req)

	response, err := parseResponse(t, w)

	if err != nil {
		t.Errorf("couldn't get response body!")
		return
	}

	if response["error"] != "Couldn't convert userId to uint!" {
		t.Errorf("expected error 'Couldn't convert userId to uint!', got %d", response["error"])
		return
	}

	// Test correct upvote
	vote := &models.Vote{
		ID:     voteId,
		PostId: postId,
		UserId: 10,
		Vote:   1,
	}

	//var posts []*models.Post

	postsRepository.EXPECT().Get("").Return(post, nil)
	votesRepository.EXPECT().Create(vote).Return(vote, nil)
	postsRepository.EXPECT().Update(post, []primitive.E{
		{"upvote_percentage", 100},
		{"score", 1},
	}).Return(postAfterUpvote, nil)

	req = httptest.NewRequest("GET", "/api/post/"+postIdString+"/upvote", nil)
	userData.Id = "10"
	ctx = req.Context()
	ctx = context.WithValue(ctx, configs.UserCtx, userData)
	req = req.WithContext(ctx)
	w = httptest.NewRecorder()

	votesHandler.Upvote(w, req)

	response, err = parseResponse(t, w)

	if err != nil {
		t.Errorf("couldn't get response body!")
		return
	}
	fmt.Printf("response is: %v", response)
}

func parseResponse(t *testing.T, w *httptest.ResponseRecorder) (map[string]interface{}, error) {
	resp := w.Result()
	responseBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Errorf("couldn't get response body!")
		return nil, err
	}

	var response map[string]interface{}
	err = json.Unmarshal(responseBytes, &response)

	if err != nil {
		t.Errorf("couldn't unmarshal response body!")
		return nil, err
	}

	return response, nil
}

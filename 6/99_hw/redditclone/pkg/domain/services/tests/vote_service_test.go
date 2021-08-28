package tests

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"redditclone/pkg/domain/models"
	"redditclone/pkg/domain/repositories/mocks"
	"redditclone/pkg/domain/services/votes"
	"testing"
)

func TestVotesServiceApplyVote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	votesRepository := mocks.NewMockIVotesRepository(ctrl)
	postsRepository := mocks.NewMockIPostsRepository(ctrl)
	votesService := votes.NewVotesService(postsRepository, votesRepository)

	postIdString := "61059ec7adc529aef3ac11b3"
	postId, err := primitive.ObjectIDFromHex(postIdString)

	if err != nil {
		t.Errorf("couldn't get postId!")
		return
	}

	//voteIdString := "6106e9ec5b7733bc045a85c4"
	//voteId, err := primitive.ObjectIDFromHex(voteIdString)

	//if err != nil {
	//	t.Errorf("couldn't get voteId!")
	//	return
	//}

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
		Votes:            []*models.Vote{},
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
	vote := &models.Vote{
		PostId: postId,
		UserId: 10,
		Vote:   1,
	}

	postsRepository.EXPECT().Get(postIdString).Return(post, nil)
	votesRepository.EXPECT().Create(vote).
		Return(gomock.AssignableToTypeOf(vote), nil).
		SetArg(0, *vote)
	postsRepository.EXPECT().Update(post, []primitive.E{
		{"upvote_percentage", 100},
		{"score", 1},
	}).Return(postAfterUpvote, nil)

	//req := httptest.NewRequest("GET", "/api/post/14/upvote", nil)
	//ctx := req.Context()
	//ctx = context.WithValue(ctx, configs.UserCtx, userData)
	//req = req.WithContext(ctx)
	//w := httptest.NewRecorder()

	result, err := votesService.ApplyVote(postIdString, uint(10), 1)

	if err != nil {
		t.Errorf("unexpected error!")
		return
	}

	fmt.Printf("result is: %v", result)
}

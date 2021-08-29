package tests

import (
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"redditclone/pkg/domain/models"
	"redditclone/pkg/domain/repositories/mocks"
	"redditclone/pkg/domain/services/votes"
	"reflect"
	"testing"
)

/**
 * @Description: Test ApplyVote function
 * @param t
 */
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

	post := &models.Post{
		ID:               postId,
		Title:            "Some post",
		Score:            0,
		UpvotePercentage: 0,
		Votes:            []*models.Vote{},
	}
	postAfterUpvote := &models.Post{
		ID:               postId,
		Title:            "Some post",
		Score:            1,
		UpvotePercentage: 100,
		Votes:            []*models.Vote{},
	}
	vote := &models.Vote{
		PostId: postId,
		UserId: 10,
		Vote:   1,
	}

	postsRepository.EXPECT().Get(postIdString).Return(post, nil)
	votesRepository.EXPECT().Create(gomock.AssignableToTypeOf(vote)).Return(vote, nil)
	postsRepository.EXPECT().Update(post, []primitive.E{
		{"upvote_percentage", 100},
		{"score", 1},
	}).Return(postAfterUpvote, nil)

	result, err := votesService.ApplyVote(postIdString, uint(10), 1)

	if err != nil {
		t.Errorf("unexpected error!")
		return
	}

	if !reflect.DeepEqual(result, postAfterUpvote) {
		t.Errorf("posts are not equal! \n expected: %v, \n got: %v", postAfterUpvote, result)
	}
}

/**
 * @Description: Test ApplyVote function where posts.Get returns error
 * @param t
 */
func TestVotesServiceApplyVoteWithPostsRepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	votesRepository := mocks.NewMockIVotesRepository(ctrl)
	postsRepository := mocks.NewMockIPostsRepository(ctrl)
	votesService := votes.NewVotesService(postsRepository, votesRepository)

	postIdString := "InvalidId"

	postsRepository.EXPECT().Get(postIdString).Return(nil, primitive.ErrInvalidHex)

	_, err := votesService.ApplyVote(postIdString, uint(10), 1)

	if err == nil {
		t.Errorf("expected the error, got nil!")
		return
	}
}

/**
 * @Description: Test ApplyVote when post already has vote from the user
 * @param t
 */
func TestVotesServiceApplyVoteUnvote(t *testing.T) {
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

	post := &models.Post{
		ID:               postId,
		Title:            "Some post",
		Score:            1,
		UpvotePercentage: 100,
		Votes: []*models.Vote{{
			PostId: postId,
			UserId: 10,
			Vote:   1,
		}},
	}
	postAfterUpvote := &models.Post{
		ID:               postId,
		Title:            "Some post",
		Score:            0,
		UpvotePercentage: 0,
		Votes:            []*models.Vote{},
	}
	vote := &models.Vote{
		PostId: postId,
		UserId: 10,
		Vote:   1,
	}

	postsRepository.EXPECT().Get(gomock.AssignableToTypeOf(postIdString)).Return(post, nil)
	postsRepository.EXPECT().Get(gomock.AssignableToTypeOf(postIdString)).Return(post, nil)
	votesRepository.EXPECT().Delete(vote.ID).Return(true, nil)
	postsRepository.EXPECT().Update(post, []primitive.E{
		{"upvote_percentage", 0},
		{"score", 0},
	}).Return(postAfterUpvote, nil)

	result, err := votesService.ApplyVote(postIdString, uint(10), 1)

	if err != nil {
		t.Errorf("unexpected error!")
		return
	}

	if !reflect.DeepEqual(result, postAfterUpvote) {
		t.Errorf("posts are not equal! \n expected: %v, \n got: %v", postAfterUpvote, result)
	}
}

/**
 * @Description: Test 1.unvote successful 2.unvote with error 3.unvote post with empty votes
 * @param t
 */
func TestVotesServiceUnvote(t *testing.T) {
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

	post := &models.Post{
		ID:               postId,
		Title:            "Some post",
		Score:            1,
		UpvotePercentage: 100,
		Votes: []*models.Vote{{
			PostId: postId,
			UserId: 10,
			Vote:   1,
		}},
	}
	postAfterUnvote := &models.Post{
		ID:               postId,
		Title:            "Some post",
		Score:            0,
		UpvotePercentage: 0,
		Votes:            []*models.Vote{},
	}
	vote := &models.Vote{
		PostId: postId,
		UserId: 10,
		Vote:   1,
	}

	// unvote successful
	postsRepository.EXPECT().Get(postIdString).Return(post, nil)
	votesRepository.EXPECT().Delete(vote.ID).Return(true, nil)
	postsRepository.EXPECT().Update(post, []primitive.E{
		{"upvote_percentage", 0},
		{"score", 0},
	}).Return(postAfterUnvote, nil)

	result, err := votesService.Unvote(uint(10), postIdString)

	if err != nil {
		t.Errorf("unexpected error!")
		return
	}

	if !reflect.DeepEqual(result, postAfterUnvote) {
		t.Errorf("posts are not equal! \n expected: %v, \n got: %v", postAfterUnvote, result)
	}

	// Unvote with error (invalid post id)
	postIdString = "InvalidId"

	postsRepository.EXPECT().Get(postIdString).Return(nil, primitive.ErrInvalidHex)

	_, err = votesService.Unvote(uint(10), postIdString)

	if err == nil {
		t.Errorf("expected the error, got nil!")
		return
	}

	// Unvote post with no votes
	post.Votes = []*models.Vote{}
	post.Score = 0
	post.UpvotePercentage = 0

	postsRepository.EXPECT().Get(postIdString).Return(post, nil)
	postsRepository.EXPECT().Update(post, []primitive.E{
		{"upvote_percentage", 0},
		{"score", 0},
	}).Return(postAfterUnvote, nil)

	result, err = votesService.Unvote(uint(10), postIdString)

	if err != nil {
		t.Errorf("unexpected error!")
		return
	}

	if !reflect.DeepEqual(result, postAfterUnvote) {
		t.Errorf("posts are not equal! \n expected: %v, \n got: %v", postAfterUnvote, result)
	}
}

/**
 * @Description: Test CalculateScore
 * @param t
 */
func TestVotesServiceCalculateScore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	votesRepository := mocks.NewMockIVotesRepository(ctrl)
	postsRepository := mocks.NewMockIPostsRepository(ctrl)
	votesService := votes.NewVotesService(postsRepository, votesRepository)

	type testCase struct {
		Votes  []*models.Vote
		Result int
	}
	testData := []testCase{
		{
			Votes:  []*models.Vote{{Vote: 1}, {Vote: 1}, {Vote: 1}},
			Result: 3,
		},{
			Votes:  []*models.Vote{{Vote: 1}, {Vote: -1}, {Vote: 2}},
			Result: 2,
		},{
			Votes:  []*models.Vote{{Vote: 1}, {Vote: -2}, {Vote: 1}},
			Result: 0,
		},
	}

	for _, tcase := range testData {
		result := votesService.CalculateScore(tcase.Votes)

		if result != tcase.Result {
			t.Errorf("wrong result! \n expected: %v \n got: %v \n", tcase.Result, result)
			return
		}
	}
}

/**
 * @Description: Test CalculateUpvotePercentage
 * @param t
 */
func TestVotesServiceCalculateUpvotePercentage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	votesRepository := mocks.NewMockIVotesRepository(ctrl)
	postsRepository := mocks.NewMockIPostsRepository(ctrl)
	votesService := votes.NewVotesService(postsRepository, votesRepository)

	type testCase struct {
		Votes  []*models.Vote
		Result int
	}
	testData := []testCase{
		{
			Votes:  []*models.Vote{{Vote: 1}, {Vote: 1}, {Vote: 1}},
			Result: 100,
		},{
			Votes:  []*models.Vote{{Vote: 1}, {Vote: -1}, {Vote: 1}},
			Result: 66,
		},{
			Votes:  []*models.Vote{{Vote: -1}, {Vote: -1}, {Vote: -1}},
			Result: 0,
		},{
			Votes:  []*models.Vote{{Vote: -1}, {Vote: -1}, {Vote: -1}, {Vote: 1}},
			Result: 25,
		},
	}

	for _, tcase := range testData {
		result := votesService.CalculateUpvotePercentage(tcase.Votes)

		if result != tcase.Result {
			t.Errorf("wrong result! \n expected: %v \n got: %v \n", tcase.Result, result)
			return
		}
	}
}

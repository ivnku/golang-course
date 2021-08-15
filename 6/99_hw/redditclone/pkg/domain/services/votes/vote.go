package votes

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"redditclone/pkg/domain/models"
	"redditclone/pkg/domain/repositories"
)

/**
 * @Description: Calculate upvote percentage
 * @param votes
 * @return int
 */
func CalculateUpvotePercentage(votes []*models.Vote) int {
	totalVotes := len(votes)
	upvotes := 0
	for _, vote := range votes {
		if vote.Vote == 1 {
			upvotes++
		}
	}

	var upvotePercentage int
	if totalVotes > 0 {
		upvotePercentage = int(float32(upvotes) / float32(totalVotes) * 100)

	}

	return upvotePercentage
}

/**
 * @Description: Calculate score
 * @param votes
 * @return int
 */
func CalculateScore(votes []*models.Vote) int {
	score := 0
	for _, vote := range votes {
		score += vote.Vote
	}

	return score
}

/**
 * @Description: Apply either upvote, downvote or unvote
 * @param postsRepository
 * @param votesRepository
 * @param postId
 * @param userId
 * @param voteValue
 * @return *models.Post
 * @return error
 */
func ApplyVote(
	postsRepository repositories.IPostsRepository,
	votesRepository repositories.IVotesRepository,
	postId string,
	userId uint,
	voteValue int,
) (*models.Post, error) {
	post, err := postsRepository.Get(postId)

	if err != nil {
		return nil, err
	}

	isAbleToVote := true

	for _, vote := range post.Votes {
		if vote.UserId == userId {
			isAbleToVote = false
			break
		}
	}

	if isAbleToVote {
		vote := &models.Vote{
			ID:     primitive.NewObjectID(),
			PostId: post.ID,
			UserId: userId,
			Vote:   voteValue,
		}

		vote, err = votesRepository.Create(vote)

		if err != nil {
			return post, err
		}

		post.Votes = append(post.Votes, vote)
		post.UpvotePercentage = CalculateUpvotePercentage(post.Votes)
		post.Score = CalculateScore(post.Votes)

		post, err = postsRepository.Update(post, []primitive.E{
			{"upvote_percentage", post.UpvotePercentage},
			{"score", post.Score},
		})

		if err != nil {
			return post, err
		}
	} else {
		post, err = Unvote(postsRepository, votesRepository, userId, post.ID.String())
	}

	return post, nil
}

/**
 * @Description: Remove user's vote from the post
 * @param votesRepository
 * @param vote
 * @param post
 * @return *models.Post
 * @return error
 */
func Unvote(
	postsRepository repositories.IPostsRepository,
	votesRepository repositories.IVotesRepository,
	userId uint,
	postId string,
) (*models.Post, error) {

	var voteId primitive.ObjectID
	post, err := postsRepository.Get(postId)

	if err != nil {
		return nil, err
	}

	for _, vote := range post.Votes {
		if vote.UserId == userId {
			voteId = vote.ID
			_, err := votesRepository.Delete(vote.ID)
			if err != nil {
				return nil, err
			}
			break
		}
	}

	for index, postVote := range post.Votes {
		if postVote.ID == voteId {
			post.Votes = append(post.Votes[:index], post.Votes[index+1:]...)
			break
		}
	}

	post.Score = CalculateScore(post.Votes)
	post.UpvotePercentage = CalculateUpvotePercentage(post.Votes)

	post, err = postsRepository.Update(post, []primitive.E{
		{"upvote_percentage", post.UpvotePercentage},
		{"score", post.Score},
	})

	if err != nil {
		return nil, err
	}

	return post, nil
}

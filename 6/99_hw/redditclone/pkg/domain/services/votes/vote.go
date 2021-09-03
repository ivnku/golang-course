package votes

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"redditclone/pkg/domain/models"
	"redditclone/pkg/domain/repositories"
)

type IVotesService interface {
	CalculateUpvotePercentage(votes []*models.Vote) int
	CalculateScore(votes []*models.Vote) int
	ApplyVote(postId string, userId uint, voteValue int) (*models.Post, error)
	Unvote(userId uint, postId string) (*models.Post, error)
}

type VotesService struct {
	postsRepository repositories.IPostsRepository
	votesRepository repositories.IVotesRepository
}

func NewVotesService(postsRepository repositories.IPostsRepository, votesRepository repositories.IVotesRepository) *VotesService {
	return &VotesService{
		postsRepository: postsRepository,
		votesRepository: votesRepository,
	}
}

/**
 * @Description: Calculate upvote percentage
 * @param votes
 * @return int
 */
func (vs *VotesService) CalculateUpvotePercentage(votes []*models.Vote) int {
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
func (vs *VotesService) CalculateScore(votes []*models.Vote) int {
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
func (vs *VotesService) ApplyVote(postId string, userId uint, voteValue int) (*models.Post, error) {
	post, err := vs.postsRepository.Get(postId)

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

		vote, err = vs.votesRepository.Create(vote)

		if err != nil {
			return post, err
		}

		post.Votes = append(post.Votes, vote)
		post.UpvotePercentage = vs.CalculateUpvotePercentage(post.Votes)
		post.Score = vs.CalculateScore(post.Votes)

		post, err = vs.postsRepository.Update(post, []primitive.E{
			{"upvote_percentage", post.UpvotePercentage},
			{"score", post.Score},
		})

		if err != nil {
			return post, err
		}
	} else {
		post, err = vs.Unvote(userId, post.ID.String())
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
func (vs *VotesService) Unvote(userId uint, postId string) (*models.Post, error) {

	var voteId primitive.ObjectID
	post, err := vs.postsRepository.Get(postId)

	if err != nil {
		return nil, err
	}

	for _, vote := range post.Votes {
		if vote.UserId == userId {
			voteId = vote.ID
			_, err := vs.votesRepository.Delete(vote.ID)
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

	post.Score = vs.CalculateScore(post.Votes)
	post.UpvotePercentage = vs.CalculateUpvotePercentage(post.Votes)

	post, err = vs.postsRepository.Update(post, []primitive.E{
		{"upvote_percentage", post.UpvotePercentage},
		{"score", post.Score},
	})

	if err != nil {
		return nil, err
	}

	return post, nil
}

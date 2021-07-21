package votes

import "redditclone/pkg/domain/models"

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

	upvotePercentage := (upvotes / totalVotes) * 100

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
package repositories

import (
	"gorm.io/gorm"
	"redditclone/pkg/domain/models"
)

type VotesRepository struct {
	db *gorm.DB
}

func NewVotesRepository(db *gorm.DB) VotesRepository {
	return VotesRepository{db}
}

/**
 * @Description: Create a new vote
 * @receiver r
 * @param vote
 * @return *Vote
 * @return error
 */
func (r *VotesRepository) Create(vote *models.Vote) (*models.Vote, error) {
	db := r.db.Create(vote)

	if err := db.Error; err != nil {
		return nil, err
	}

	return vote, nil
}

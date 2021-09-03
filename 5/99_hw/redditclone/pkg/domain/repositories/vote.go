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

/**
 * @Description: Delete the vote
 * @receiver r *PostsRepository
 * @return bool
 * @return error
 */
func (r *VotesRepository) Delete(id uint) (bool, error) {

	db := r.db.Delete(&models.Vote{}, id)

	if err := db.Error; err != nil {
		return false, err
	}

	return true, nil
}
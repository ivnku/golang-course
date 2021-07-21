package repositories

import (
	"gorm.io/gorm"
	"redditclone/pkg/domain/models"
)

type UsersRepository struct {
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) UsersRepository {
	return UsersRepository{db}
}

/**
 * @Description: Get the list of users
 * @receiver r
 * @return []*User
 * @return error
 */
func (r *UsersRepository) List() ([]*models.User, error) {
	var users []*models.User

	db := r.db.Find(&users)

	if err := db.Error; err != nil {
		return users, err
	}

	return users, nil
}

/**
 * @Description: Get a user by id
 * @receiver r
 * @param id
 * @return *User
 * @return error
 */
func (r *UsersRepository) Get(id int) (*models.User, error) {
	var user *models.User

	db := r.db.First(&user, id)

	if err := db.Error; err != nil {
		return nil, err
	}

	return user, nil
}

/**
 * @Description: Get a user by his name
 * @receiver r
 * @param name
 * @return *User
 * @return error
 */
func (r *UsersRepository) GetByName(name string) (*models.User, error) {
	var user *models.User

	db := r.db.First(&user, "name = ?", name)

	if err := db.Error; err != nil {
		return nil, err
	}

	return user, nil
}

/**
 * @Description: Create a user row in the db
 * @receiver r
 * @param user
 * @return uint
 * @return error
 */
func (r *UsersRepository) Create(user *models.User) (uint, error) {
	result := r.db.Create(&user)

	if result.Error != nil {
		return 0, result.Error
	}

	return user.ID, nil
}

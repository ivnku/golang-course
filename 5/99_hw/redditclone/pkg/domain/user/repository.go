package user

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{db}
}

/**
 * @Description: Get the list of users
 * @receiver r
 * @return []*User
 * @return error
 */
func (r *Repository) List() ([]*User, error) {
	var users []*User

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
func (r *Repository) Get(id int) (*User, error) {
	var user *User

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
func (r *Repository) GetByName(name string) (*User, error) {
	var user *User

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
func (r *Repository) Create(user *User) (uint, error) {
	result := r.db.Create(&user)

	if result.Error != nil {
		return 0, result.Error
	}

	return user.ID, nil
}

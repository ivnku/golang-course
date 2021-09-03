package repositories

import (
	"gorm.io/gorm"
	"redditclone/pkg/domain/models"
)

type PostsRepository struct {
	db *gorm.DB
}

func NewPostsRepository(db *gorm.DB) PostsRepository {
	return PostsRepository{db}
}

/**
 * @Description: Create a new post
 * @receiver r
 * @param post
 * @return *Post
 * @return error
 */
func (r *PostsRepository) Create(post *models.Post) (*models.Post, error) {
	db := r.db.Create(post)

	if err := db.Error; err != nil {
		return nil, err
	}

	return post, nil
}

/**
 * @Description: Update the post
 * @receiver r
 * @param post
 * @return *Post
 * @return error
 */
func (r *PostsRepository) Update(post *models.Post, fields []string) (*models.Post, error) {
	var db *gorm.DB
	if fields != nil {
		db = r.db.Model(post).Select(fields).Updates(post)
	} else {
		db = r.db.Save(post)
	}

	if err := db.Error; err != nil {
		return post, err
	}

	return post, nil
}

/**
 * @Description: Get one post by id
 * @receiver r *PostsRepository
 * @param id uint
 * @return *Post
 * @return error
 */
func (r *PostsRepository) Get(id uint) (*models.Post, error) {
	var post *models.Post

	db := r.db.Preload("Comments.User").Preload("Votes").Joins("User").First(&post, id)

	if err := db.Error; err != nil {
		return post, err
	}

	return post, nil
}

/**
 * @Description: Get the list of posts
 * @receiver r *PostsRepository
 * @return []*Post
 * @return error
 */
func (r *PostsRepository) List() ([]*models.Post, error) {
	var posts []*models.Post

	db := r.db.Preload("Comments.User").Preload("Votes").Joins("User").Find(&posts)

	if err := db.Error; err != nil {
		return posts, err
	}

	return posts, nil
}

/**
 * @Description: Get posts within a certain category
 * @receiver r
 * @param categoryName
 * @return []*Post
 * @return error
 */
func (r *PostsRepository) CategoryList(categoryName string) ([]*models.Post, error) {
	var posts []*models.Post

	db := r.db.Preload("Comments.User").Preload("Votes").Joins("User").Find(&posts, "category = ?", categoryName)

	if err := db.Error; err != nil {
		return posts, err
	}

	return posts, nil
}

/**
 * @Description: Get posts of a certain user
 * @receiver r
 * @param categoryName
 * @return []*Post
 * @return error
 */
func (r *PostsRepository) UserList(userId uint) ([]*models.Post, error) {
	var posts []*models.Post

	db := r.db.Preload("Comments.User").Preload("Votes").Joins("User").Find(&posts, "user_id = ?", userId)

	if err := db.Error; err != nil {
		return posts, err
	}

	return posts, nil
}

/**
 * @Description: Delete the post
 * @receiver r *PostsRepository
 * @return bool
 * @return error
 */
func (r *PostsRepository) Delete(id uint) (bool, error) {

	db := r.db.Delete(&models.Post{}, id)

	if err := db.Error; err != nil {
		return false, err
	}

	return true, nil
}

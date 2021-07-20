package post

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
 * @Description: Create a new post
 * @receiver r
 * @param post
 * @return *Post
 * @return error
 */
func (r *Repository) Create(post *Post) (*Post, error) {
	db := r.db.Create(post)

	if err := db.Error; err != nil {
		return nil, err
	}

	return post, nil
}

/**
 * @Description: Get one post by id
 * @receiver r *Repository
 * @param id uint
 * @return *Post
 * @return error
 */
func (r *Repository) Get(id uint) (*Post, error) {
	var post *Post

	db := r.db.Preload("Comments.User").Preload("Votes").Joins("User").First(&post, id)

	if err := db.Error; err != nil {
		return post, err
	}

	return post, nil
}

/**
 * @Description: Get the list of posts
 * @receiver r *Repository
 * @return []*Post
 * @return error
 */
func (r *Repository) List() ([]*Post, error) {
	var posts []*Post

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
func (r *Repository) CategoryList(categoryName string) ([]*Post, error) {
	var posts []*Post

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
func (r *Repository) UserList(userId uint) ([]*Post, error) {
	var posts []*Post

	db := r.db.Preload("Comments.User").Preload("Votes").Joins("User").Find(&posts, "user_id = ?", userId)

	if err := db.Error; err != nil {
		return posts, err
	}

	return posts, nil
}

/**
 * @Description: Delete the post
 * @receiver r *Repository
 * @return bool
 * @return error
 */
func (r *Repository) Delete(id uint) (bool, error) {

	db := r.db.Delete(&Post{}, id)

	if err := db.Error; err != nil {
		return false, err
	}

	return true, nil
}

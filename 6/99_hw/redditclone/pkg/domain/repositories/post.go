package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"redditclone/pkg/domain/models"
)

type PostsRepository struct {
	db         *mongo.Client
	collection *mongo.Collection
}

func NewPostsRepository(db *mongo.Client) PostsRepository {
	collection := db.Database("redditclone").Collection("posts")
	return PostsRepository{db, collection}
}

/**
 * @Description: Create a new post
 * @receiver r
 * @param post
 * @return *Post
 * @return error
 */
func (r *PostsRepository) Create(post *models.Post) (*models.Post, error) {
	var ctx = context.Background()
	post.ID = primitive.NewObjectID()
	result, err := r.collection.InsertOne(ctx, post)

	if err != nil {
		return nil, err
	}

	post.ID = result.InsertedID.(primitive.ObjectID)

	return post, nil
}

/**
 * @Description: Update the post
 * @receiver r
 * @param post
 * @return *Post
 * @return error
 */
func (r *PostsRepository) Update(post *models.Post, fields []primitive.E) (*models.Post, error) {
	var ctx = context.Background()
	fieldsToUpdate := bson.D{}
	for _, field := range fields {
		fieldsToUpdate = append(fieldsToUpdate, primitive.E{Key: "$set", Value: bson.D{field}})
	}

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": post.ID},
		fieldsToUpdate,
	)

	if err != nil {
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
func (r *PostsRepository) Get(id string) (*models.Post, error) {
	var ctx = context.Background()
	postId, err := primitive.ObjectIDFromHex(id)

	qry := []bson.M{
		{
			"$match": bson.M{
				"_id": postId,
			},
		},
		{"$lookup": bson.M{
				"from":         "comments", // Child collection to join
				"localField":   "_id",      // Parent collection reference holding field
				"foreignField": "post_id",  // Child collection reference field
				"as":           "comments", // Arbitrary field name to store result set
		}},
		{"$lookup": bson.M{
			"from":         "votes",
			"localField":   "_id",
			"foreignField": "post_id",
			"as":           "votes",
		}},
	}

	if err != nil {
		return nil, err
	}

	cur, err := r.collection.Aggregate(ctx, qry)

	if err != nil {
		return nil, err
	}

	var posts []*models.Post
	if err := cur.All(context.Background(), &posts); err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	return posts[0], nil
}

/**
 * @Description: Get the list of posts
 * @receiver r *PostsRepository
 * @return []*Post
 * @return error
 */
func (r *PostsRepository) List() ([]*models.Post, error) {
	var ctx = context.Background()
	var posts []*models.Post

	qry := []bson.M{
		{"$lookup": bson.M{
			"from":         "comments", // Child collection to join
			"localField":   "_id",      // Parent collection reference holding field
			"foreignField": "post_id",  // Child collection reference field
			"as":           "comments", // Arbitrary field name to store result set
		}},
		{"$lookup": bson.M{
			"from":         "votes",
			"localField":   "_id",
			"foreignField": "post_id",
			"as":           "votes",
		}},
	}

	cur, err := r.collection.Aggregate(ctx, qry)

	if err != nil {
		return posts, err
	}

	if err := cur.All(context.Background(), &posts); err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

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
	var ctx = context.Background()
	var posts []*models.Post

	qry := []bson.M{
		{"$match": bson.M{"category": categoryName}},
		{"$lookup": bson.M{
			"from":         "comments",
			"localField":   "_id",
			"foreignField": "post_id",
			"as":           "comments",
		},
		},
		{"$lookup": bson.M{
			"from":         "votes",
			"localField":   "_id",
			"foreignField": "post_id",
			"as":           "votes",
		},
		},
	}

	cur, err := r.collection.Aggregate(ctx, qry)

	if err != nil {
		return posts, err
	}

	if err := cur.All(context.Background(), &posts); err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

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
	var ctx = context.Background()
	var posts []*models.Post

	qry := []bson.M{
		{"$match": bson.M{"user.id": userId}},
		{"$lookup": bson.M{
			"from":         "comments",
			"localField":   "_id",
			"foreignField": "post_id",
			"as":           "comments",
		}},
		{"$lookup": bson.M{
			"from":         "votes",
			"localField":   "_id",
			"foreignField": "post_id",
			"as":           "votes",
		}},
	}

	cur, err := r.collection.Aggregate(ctx, qry)

	if err != nil {
		return posts, err
	}

	if err := cur.All(context.Background(), &posts); err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	return posts, nil
}

/**
 * @Description: Delete the post
 * @receiver r *PostsRepository
 * @return bool
 * @return error
 */
func (r *PostsRepository) Delete(id string) (bool, error) {
	var ctx = context.Background()
	postId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return false, err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": postId})

	if err != nil {
		return false, err
	}

	return true, nil
}

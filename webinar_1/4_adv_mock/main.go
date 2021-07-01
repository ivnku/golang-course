package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostsRepo struct {
	Collection IMongoCollection
}

type Post struct {
	ID    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Email string             `json:"email" bson:"email"`
}

func (repo *PostsRepo) GetOne(ctx context.Context, objID primitive.ObjectID) (*Post, error) {
	post := &Post{}
	err := repo.Collection.FindOne(ctx, bson.M{"_id": objID}).Decode(post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	objID := primitive.NewObjectID()

	mongoDatabase := client.Database("adb_mock")
	coll := mongoDatabase.Collection("posts")
	postsCollection := &MongoCollection{
		Сoll: coll,
	}

	coll.InsertOne(ctx, &Post{
		ID:    objID,
		Email: "romanov.vasily@mail.ru",
	})

	repo := &PostsRepo{
		Collection: postsCollection,
	}

	res, err := repo.GetOne(ctx, objID)
	fmt.Printf("result %#v, err: %v", res, err)
}

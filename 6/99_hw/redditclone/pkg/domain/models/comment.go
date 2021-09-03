package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Comment struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	PostID  primitive.ObjectID `json:"-" bson:"post_id"`
	UserID  uint               `json:"-"`
	User    User               `json:"author" bson:"author"`
	Body    string             `json:"body" bson:"body"`
	Created time.Time          `json:"created" bson:"created"`
}

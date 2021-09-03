package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Vote struct {
	ID     primitive.ObjectID `json:"-" bson:"_id"`
	PostId primitive.ObjectID `json:"-" bson:"post_id"`
	UserId uint               `json:"user,string" bson:"user"`
	Vote   int                `json:"vote" bson:"vote"`
}

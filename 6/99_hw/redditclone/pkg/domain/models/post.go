package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Post struct {
	ID               primitive.ObjectID `json:"id" bson:"_id"`
	Title            string             `json:"title" bson:"title"`
	Category         string             `json:"category" bson:"category"`
	Score            int                `json:"score" bson:"score"`
	Type             string             `json:"type" bson:"type"`
	Url              string             `json:"url" bson:"url"`
	Text             string             `json:"text" bson:"text"`
	UpvotePercentage int                `json:"upvotePercentage" bson:"upvotePercentage"`
	Views            int                `json:"views" bson:"views"`
	CreatedAt        time.Time          `json:"created" bson:"created"`
	User             User               `json:"author"`
	Comments         []*Comment         `json:"comments"`
	Votes            []*Vote            `json:"votes"`
}

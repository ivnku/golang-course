package post

import (
	"redditclone/pkg/domain/comment"
	"redditclone/pkg/domain/user"
	"redditclone/pkg/domain/vote"
)

type Post struct {
	ID               uint              `json:"id"`
	Title            string            `json:"title"`
	Category         string            `json:"category"`
	Score            int               `json:"score"`
	Type             string            `json:"type"`
	Url              string            `json:"url"`
	Text             string            `json:"text"`
	UpvotePercentage int               `json:"upvotePercentage"`
	Views            int               `json:"views"`
	CreatedAt        string            `json:"created"`
	UserID           uint              `json:"-"`
	User             user.User         `json:"author"`
	Comments         []comment.Comment `json:"comments"`
	Votes            []*vote.Vote      `json:"votes"`
}

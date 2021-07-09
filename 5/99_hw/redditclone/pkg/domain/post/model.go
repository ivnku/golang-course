package post

import (
	"redditclone/pkg/domain/comment"
	"redditclone/pkg/domain/user"
	"redditclone/pkg/domain/vote"
)

type Post struct {
	ID   uint
	Title string
	Category string
	Score int
	Type string
	UpvotePercentage int
	Views int
	CreatedAt string
	Author *user.User
	Comments []*comment.Comment
	Votes []*vote.Vote
}
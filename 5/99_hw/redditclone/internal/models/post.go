package models

type Post struct {
	ID   uint
	Title string
	Category string
	Score int
	Type string
	UpvotePercentage int
	Views int
	CreatedAt string
	Author *User
	Comments []*Comment
	Votes []*Vote
}
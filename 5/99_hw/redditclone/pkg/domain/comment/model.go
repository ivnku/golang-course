package comment

import (
	"redditclone/pkg/domain/user"
)

type Comment struct {
	ID      uint      `json:"id"`
	PostID  uint      `json:"-"`
	UserID  uint      `json:"-"`
	User    user.User `json:"author"`
	Body    string    `json:"body"`
	Created string    `json:"created"`
}

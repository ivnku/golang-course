package models

type Comment struct {
	ID      uint   `json:"id"`
	PostID  uint   `json:"-"`
	UserID  uint   `json:"-"`
	User    User   `json:"author"`
	Body    string `json:"body"`
	Created string `json:"created"`
}

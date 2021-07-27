package models

type Comment struct {
	ID      uint   `json:"id"`
	PostID  string `json:"-"`
	UserID  uint   `json:"-"`
	User    User   `json:"author"`
	Body    string `json:"body"`
	Created string `json:"created"`
}

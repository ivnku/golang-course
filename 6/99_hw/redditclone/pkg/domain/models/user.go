package models

type User struct {
	ID       uint   `json:"id,string"`
	Name     string `json:"username"`
	Password string `json:"-"`
}

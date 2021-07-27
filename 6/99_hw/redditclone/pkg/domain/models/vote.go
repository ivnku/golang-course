package models

type Vote struct {
	ID     uint   `json:"-"`
	PostId string `json:"-"`
	UserId uint   `json:"user"`
	Vote   int    `json:"vote"`
}

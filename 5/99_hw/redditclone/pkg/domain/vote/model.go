package vote

type Vote struct {
	ID     uint `json:"-"`
	PostId uint `json:"-"`
	UserId uint `json:"user"`
	Vote   int  `json:"vote"`
}

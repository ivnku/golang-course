package user

type User struct {
	ID   uint   `json:"id"`
	Name string `json:"username"`
	Password string `json:"-"`
}

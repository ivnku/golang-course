package models

type Post struct {
	ID               uint              `json:"id"`
	Title            string       `json:"title"`
	Category         string       `json:"category"`
	Score            int       `json:"score"`
	Type             string    `json:"type"`
	Url              string    `json:"url"`
	Text             string    `json:"text"`
	UpvotePercentage int       `json:"upvotePercentage"`
	Views            int       `json:"views"`
	CreatedAt        string    `json:"created"`
	UserID           uint      `json:"-"`
	User             User      `json:"author"`
	Comments         []Comment `json:"comments"`
	Votes            []*Vote   `json:"votes"`
}

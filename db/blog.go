package db

type Blog struct {
	Blogid string `json:"blogid"`
	BlogTitle string `json:"blog_title"`
	BlogContent string `json:"blog_content"`
	UserId string `json:"user_id"`
	Like string `json:"like"`
	CreateAt int64 `json:"create_at"`
	CommentId string `json:"comment_id"`
	Pictures []string `json:"pictures"`
}

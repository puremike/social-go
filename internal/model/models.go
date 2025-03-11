package model

type PostModel struct {
	ID int `json:"id"`
	Content string	`json:"content"`
	Title string 	`json:"title"`
	UserID int	`json:"user_id"`
	Tags []string	`json:"tags"`
	CreatedAt string	`json:"created_at"`
	UpdatedAt string	`json:"updated_at"`
	Comments []CommentModel `json:"comments"`
}

type FollowerModel struct {
	FollowerID int `json:"follower_id"`
	UserID int `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

type UserModel struct {
	ID int `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"-"`
	CreatedAt string `json:"created_at"`
}


type CommentModel struct {
	ID        int64  `json:"id"`
	PostID    int64  `json:"post_id"`
	UserID    int64  `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	User UserModel `json:"user"`
}
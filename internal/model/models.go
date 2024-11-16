package model

type PostModel struct {
	ID int `json:"id"`
	Content string	`json:"content"`
	Title string 	`json:"title"`
	UserID int	`json:"user_id"`
	Tags []string	`json:"tags"`
	CreatedAt string	`json:"created_at"`
	UpdatedAt string	`json:"updated_at"`
}

type UserModel struct {
	ID int `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"-"`
	CreatedAt string `json:"created_at"`
}

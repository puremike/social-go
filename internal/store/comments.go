package store

import (
	"context"
	"database/sql"

	"github.com/puremike/social-go/internal/model"
)


type CommentStore struct {
	db *sql.DB
}

func (s *CommentStore) GetCommentsByPostID(ctx context.Context, id int) ([]model.CommentModel, error) {

	query := `SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, users.id, users.username FROM comments c
	JOIN users on users.id = c.user_id 
	WHERE c.post_id = $1 
	ORDER BY c.created_at DESC`

	rows, err := s.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := []model.CommentModel{}
	for rows.Next() {
		var c model.CommentModel
		c.User = model.UserModel{}

		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.User.ID, &c.User.Username)

		if err != nil {
			return nil, err
		}

		comments = append(comments, c)
	}

	return comments, nil
}

func (s *CommentStore) Create(ctx context.Context, comment *model.CommentModel) error {
	query := `
		INSERT INTO comments (post_id, user_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		comment.PostID,
		comment.UserID,
		comment.Content,
	).Scan(
		&comment.ID,
		&comment.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

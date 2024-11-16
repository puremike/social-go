package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/puremike/social-go/internal/model"
)

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *model.PostModel) error {
	query := `INSERT INTO posts (content, title, user_id, tags) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	err := s.db.QueryRowContext(ctx, query, post.Content, post.Title, post.UserID, pq.Array(post.Tags)).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return err
	}
	
	return nil
}

func (s *PostStore) GetPostByID(ctx context.Context, id int) (*model.PostModel, error) {
	query := `SELECT id, title, user_id, tags, created_at, updated_at FROM posts WHERE id = $1`
	post := &model.PostModel{}

	err := s.db.QueryRowContext(ctx, query, id).Scan(&post.ID, &post.Content, &post.UserID, pq.Array(&post.Tags), &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	return post, nil
}
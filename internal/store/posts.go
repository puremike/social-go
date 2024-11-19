package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

func (s *PostStore) GetAllPosts(ctx context.Context) ([]model.PostModel, error) {
	query := `SELECT id, content, title, user_id, tags, created_at, updated_at FROM posts;`

	posts := []model.PostModel{}

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p model.PostModel

        err := rows.Scan(&p.ID, &p.Content, &p.Title, &p.UserID, pq.Array(&p.Tags), &p.CreatedAt, &p.UpdatedAt)

        if err!= nil {
            return nil, err
        }

        posts = append(posts, p)
	}
	return posts, nil

}

func (s *PostStore) GetPostByID(ctx context.Context, id int) (*model.PostModel, error) {
	query := `SELECT id, content, title, user_id, tags, created_at, updated_at FROM posts WHERE id = $1`
	post := &model.PostModel{}

	err := s.db.QueryRowContext(ctx, query, id).Scan(&post.ID, &post.Content, &post.Title, &post.UserID, pq.Array(&post.Tags), &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	return post, nil

}

func (s *PostStore) DeletePostByID(ctx context.Context, id int) (string, error) {
	query := `DELETE FROM posts WHERE id = $1 RETURNING id, title`

	post := &model.PostModel{}
	
	err := s.db.QueryRowContext(ctx, query, id).Scan(&post.ID, &post.Title)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
            return "", ErrPostNotFound
        }
        return "", err
	}

	message := fmt.Sprintf("The post with the id (%d) and title (%s) has been deleted successfully", post.ID, post.Title)

	return message, nil 
}

func (s *PostStore) DeleteAllPosts(ctx context.Context) (string, error) {
	query := `DELETE FROM posts`
	_, err := s.db.ExecContext(ctx, query)
	
    if err != nil {
		return "", errors.New("unable to delete posts")
    }	
	message := "All posts have been deleted successfully"
	
	return message, nil
}

func (s *PostStore) UpdatePost(ctx context.Context, id int, post *model.PostModel) error {
	query := `UPDATE posts SET title = $2, content = $3, tags = $4 WHERE id = $1 RETURNING id, title, content, user_id, tags, created_at, updated_at;`
	
	err := s.db.QueryRowContext(ctx, query, id, post.Title, post.Content, pq.Array(post.Tags)).Scan(&post.ID, &post.Title, &post.Content, &post.UserID, pq.Array(&post.Tags), &post.CreatedAt, &post.UpdatedAt)

	if err !=  nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("the post with the ID (%d )does not exist", post.ID)
		}
		return err
	}

	return nil
}
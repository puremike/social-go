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

type PostWithMetaData struct {
	model.PostModel
	CommentsCount int `json:"comments_count"`
}

func (s *PostStore) Create(ctx context.Context, post *model.PostModel) error {
	query := `INSERT INTO posts (content, title, user_id, tags) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, post.Content, post.Title, post.UserID, pq.Array(post.Tags)).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostStore) GetAllPosts(ctx context.Context) ([]model.PostModel, error) {
	query := `SELECT id, content, title, user_id, tags, created_at, updated_at FROM posts;`

	posts := []model.PostModel{}

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p model.PostModel

		err := rows.Scan(&p.ID, &p.Content, &p.Title, &p.UserID, pq.Array(&p.Tags), &p.CreatedAt, &p.UpdatedAt)

		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}
	return posts, nil

}

func (s *PostStore) GetPostByID(ctx context.Context, id int) (*model.PostModel, error) {
	query := `SELECT id, content, title, user_id, tags, created_at, updated_at FROM posts WHERE id = $1`
	post := &model.PostModel{}

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, id).Scan(&post.ID, &post.Content, &post.Title, &post.UserID, pq.Array(&post.Tags), &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	return post, nil

}

func (s *PostStore) DeletePostByID(ctx context.Context, id int) error {
	query := `DELETE FROM posts WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, id)

	if err != nil {
		return fmt.Errorf("failed to delete post with ID %d: %w", id, err)
	}

	// check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no post found with ID %d", id)
	}

	return nil
}

// func (s *PostStore) DeleteAllPosts(ctx context.Context) error {
// 	query := `DELETE FROM posts`
// 	_, err := s.db.ExecContext(ctx, query)

//     if err != nil {
// 		return fmt.Errorf("unable to delete posts: %w", err)
//     }
// 	return nil
// }

func (s *PostStore) UpdatePost(ctx context.Context, id int, post *model.PostModel) error {
	query := `UPDATE posts SET title = $2, content = $3, tags = $4 WHERE id = $1 RETURNING id, title, content, user_id, tags, created_at, updated_at;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, id, post.Title, post.Content, pq.Array(post.Tags)).Scan(&post.ID, &post.Title, &post.Content, &post.UserID, pq.Array(&post.Tags), &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("the post with the ID (%d )does not exist", post.ID)
		}
		return err
	}

	return nil
}

func (s *PostStore) GetUserFeed(ctx context.Context, id int, fq PagQuery) ([]PostWithMetaData, error) {
	query := `
		SELECT 
			p.id, p.user_id, p.title, p.content, p.created_at, p.tags,
			u.username,
			COUNT(c.id) AS comments_count
		FROM posts p
		LEFT JOIN comments c ON c.post_id = p.id
		LEFT JOIN users u ON p.user_id = u.id
		JOIN followers f ON f.follower_id = p.user_id OR p.user_id = $1
		WHERE 
			f.user_id = $1 AND
			(p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%') AND
			(p.tags @> $5 OR $5 IS NULL)
		GROUP BY p.id, u.username
		ORDER BY p.created_at ` + fq.Sort + `
		LIMIT $2 OFFSET $3
		`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, id, fq.Limit, fq.Offset, fq.Search, pq.Array(fq.Tags))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var feed []PostWithMetaData
	for rows.Next() {
		var p PostWithMetaData
		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
			pq.Array(&p.Tags),
			&p.User.Username,
			&p.CommentsCount,
		)
		if err != nil {
			return nil, err
		}

		feed = append(feed, p)
	}

	return feed, nil
}

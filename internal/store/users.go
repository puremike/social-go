package store

import (
	"context"
	"database/sql"

	"github.com/puremike/social-go/internal/model"
)

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *model.UserModel) error {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, created_at`

	err := s.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
        return err
    }

	return nil
}
package store

import (
	"context"
	"database/sql"
	"errors"
)

type FollowerStore struct {
	db *sql.DB
}

func (s *FollowerStore) Follow(ctx context.Context, follower_id, id int) error {

	query := `INSERT INTO followers (user_id, follower_id)
    VALUES ($1, $2)
    ON CONFLICT (user_id, follower_id) DO NOTHING`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, id, follower_id)

	if err != nil {
		return errors.New("something went wrong")
	}

	return nil
}

func (s *FollowerStore) Unfollow(ctx context.Context, follower_id, id int) error {
	query := `DELETE FROM followers 
	WHERE user_id = $1 AND follower_id = $2`
	
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, id, follower_id)

	if err != nil {
		return errors.New("something went wrong")	
	}

	return nil
}
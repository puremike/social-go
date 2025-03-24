package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserStore struct {
	db *sql.DB
}

type Password struct {
	password *string
	hash     []byte
}

type UserModel struct {
	ID        int      `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  Password `json:"-"`
	CreatedAt string   `json:"created_at"`
	IsActive  bool     `json:"is_active"`
}

func (p *Password) Set(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.password = &password
	p.hash = hash

	return nil
}

func (s *UserStore) Create(ctx context.Context, user *UserModel) error {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, created_at`

	err := s.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) GetUserByID(ctx context.Context, id int) (*UserModel, error) {
	query := `SELECT id, username, email, password, created_at FROM users WHERE id = $1`

	user := &UserModel{}

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *UserStore) createUserForInvitation(ctx context.Context, tx *sql.Tx, user *UserModel) error {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, created_at`

	err := tx.QueryRowContext(ctx, query, user.Username, user.Email, user.Password.hash).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}
	return nil
}

func (s *UserStore) CreateAndInvite(ctx context.Context, user *UserModel, token string, invitationExp time.Duration) error {
	return withTX(s.db, ctx, func(tx *sql.Tx) error {
		// create user
		if err := s.createUserForInvitation(ctx, tx, user); err != nil {
			return err
		}
		// create user invitation
		if err := s.createUserInvitation(ctx, tx, token, invitationExp, user.ID); err != nil {
			return err
		}
		return nil
	})
}

func (s *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, invitationExp time.Duration, id int) error {

	query := `INSERT INTO user_invitation (token, user_id, expiry) VALUES ($1, $2, $3)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, id, time.Now().Add(invitationExp))

	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) Activate(ctx context.Context, token string) error {
	return withTX(s.db, ctx, func(tx *sql.Tx) error {
		// get the token belonging to the user
		user, err := s.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}

		//  update the user
		user.IsActive = true

		if err := s.update(ctx, tx, user); err != nil {
			return err
		}

		if err := s.deleteUserInvitations(ctx, tx, user.ID); err != nil {
			return err
		}
		return nil
	})
}

func (s *UserStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*UserModel, error) {
	query := `SELECT u.id, u.username, u.email, u.created_at, u.is_active
				FROM users u
				JOIN user_invitation ui ON u.id = ui.user_id
				WHERE ui.token = $1 AND ui.expiry > $2`

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &UserModel{}
	// you make use of the scan method for the QueryRowContext(a single row) and loop through the rows for QueryRow method
	if err := tx.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.IsActive); err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrUserNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}

func (s *UserStore) update(ctx context.Context, tx *sql.Tx, user *UserModel) error {
	query := `UPDATE users SET username = $1, email = $2, is_active = $3 WHERE id = $4`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.IsActive, user.ID)

	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) deleteUserInvitations(ctx context.Context, tx *sql.Tx, id int) error {
	query := `DELETE FROM user_invitation WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

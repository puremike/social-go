package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/puremike/social-go/internal/model"
	"golang.org/x/crypto/bcrypt"
)


type UserStore struct {
	db *sql.DB
}

type Password struct {
	password *string
	hash []byte
}

type UserModel struct {
	ID int `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password Password `json:"-"`
	CreatedAt string `json:"created_at"`
}


func(p *Password) Set(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.password = &password
	p.hash = hash

	return nil
}

func (s *UserStore) Create(ctx context.Context, user *model.UserModel) error {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, created_at`

	err := s.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
        return err
    }

	return nil
}

func (s *UserStore) GetUserByID(ctx context.Context, id int) (*model.UserModel, error) {
	query := `SELECT id, username, email, password, created_at FROM users WHERE id = $1`

	user := &model.UserModel{}

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


func (s *UserStore) CreateAndInvite(ctx context.Context, user *UserModel, token string) error {
	// transaction wrapper
	// create user
	// create user invite
	return nil
}
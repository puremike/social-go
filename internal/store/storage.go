package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/puremike/social-go/internal/model"
)

type Storage struct {
	Users interface {
		Create(context.Context, *UserModel) error
		GetUserByID(context.Context, int) (*UserModel, error)
		createUserForInvitation(context.Context, *sql.Tx, *UserModel) error
		CreateAndInvite(context.Context, *UserModel, string, time.Duration) error
		Activate(context.Context, string) error
		Delete(context.Context, int) error
		DeleteUserByID(context.Context, int) error
		GetUserByEmail(context.Context, string) (*UserModel, error)
	}

	Followers interface {
		Follow(context.Context, int, int) error
		Unfollow(context.Context, int, int) error
	}

	Posts interface {
		Create(context.Context, *model.PostModel) error
		GetPostByID(context.Context, int) (*model.PostModel, error)
		GetAllPosts(context.Context) ([]model.PostModel, error)
		DeletePostByID(context.Context, int) error
		// DeleteAllPosts(context.Context) error
		UpdatePost(context.Context, int, *model.PostModel) error
		GetUserFeed(context.Context, int, PagQuery) ([]PostWithMetaData, error)
	}

	Comments interface {
		GetCommentsByPostID(context.Context, int) ([]model.CommentModel, error)
		Create(context.Context, *model.CommentModel) error
	}

	Roles interface {
		GetByName(context.Context, string) (*Role, error)
	}
}

var (
	ErrPostNotFound      = errors.New("post not found")
	ErrUserNotFound      = errors.New("user not found")
	ErrDuplicateEmail    = errors.New("duplicate email")
	ErrDuplicateUsername = errors.New("duplicate username")
	QueryTimeoutDuration = 5 * time.Second
)

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Users:     &UserStore{db},
		Posts:     &PostStore{db},
		Comments:  &CommentStore{db},
		Followers: &FollowerStore{db},
		Roles:     &RoleStore{db},
	}
}

func withTX(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

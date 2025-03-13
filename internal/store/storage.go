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
		Create(context.Context, *model.UserModel) error
		GetUserByID(context.Context, int) (*model.UserModel, error)
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
		DeleteAllPosts(context.Context) error
		UpdatePost(context.Context, int, *model.PostModel) error
		GetUserFeed(context.Context, int) ([]PostWithMetaData, error)
	}

	Comments interface {
        GetCommentsByPostID(context.Context, int) ([]model.CommentModel, error)
		Create(context.Context, *model.CommentModel) error
    }
}

var (
	ErrPostNotFound = errors.New("post not found")
	ErrUserNotFound = errors.New("user not found")
	QueryTimeoutDuration = 5 * time.Second
)

func NewStorage(db *sql.DB) Storage {
	str := Storage {
		Users : &UserStore{db},
		Posts : &PostStore{db},
		Comments: &CommentStore{db},
		Followers: &FollowerStore{db},
	}

	return str
}


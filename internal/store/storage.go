package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/puremike/social-go/internal/model"
)

type Storage struct {
	Users interface {
		Create(context.Context, *model.UserModel) error
	}
	Posts interface {
		Create(context.Context, *model.PostModel) error
		GetPostByID(context.Context, int) (*model.PostModel, error)
	}
	Comments interface {
        GetCommentsByPostID(context.Context, int) ([]model.CommentModel, error)
    }
}

var ErrPostNotFound = errors.New("post not found")

func NewStorage(db *sql.DB) Storage {
	str := Storage {
		Users : &UserStore{db},
		Posts : &PostStore{db},
		Comments: &CommentStore{db},
	}

	return str
}
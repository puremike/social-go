package store

import (
	"context"
	"database/sql"

	"github.com/puremike/social-go/internal/model"
)

type Storage struct {
	Users interface {
		Create(context.Context, *model.UserModel) error
	}
	Posts interface {
		Create(context.Context, *model.PostModel) error
	}
}

func NewStorage(db *sql.DB) Storage {
	str := Storage {
		Users : &UserStore{db},
		Posts : &PostStore{db},
	}

	return str
}
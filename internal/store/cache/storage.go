package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/puremike/social-go/internal/store"
)

type Storage struct {
	Users interface {
		Get(context.Context, int) (*store.UserModel, error)
		Set(context.Context, *store.UserModel) error
		Delete(context.Context, int)
	}
}

func NewRdbStorage(rdb *redis.Client) Storage {
	return Storage{
		Users: &UserStoreRdb{rdb},
	}
}

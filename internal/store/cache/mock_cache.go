package cache

import (
	"context"

	"github.com/puremike/social-go/internal/store"
)

type mockCacheStore struct {
}

func NewMockCacheStore() Storage {
	return Storage{
		Users: &mockCacheStore{},
	}
}

func (mc *mockCacheStore) Get(ctx context.Context, id int) (*store.UserModel, error) {
	return &store.UserModel{ID: id}, nil
}
func (mc *mockCacheStore) Set(ctx context.Context, user *store.UserModel) error {
	return nil
}

func (mc *mockCacheStore) Delete(ctx context.Context, id int) {}

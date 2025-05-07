package store

import (
	"context"
	"database/sql"
	"time"
)

func NewMockStore() Storage {
	return Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct {
}

func (m *MockUserStore) Create(ctx context.Context, u *UserModel) error {
	return nil
}

func (m *MockUserStore) GetUserByID(ctx context.Context, id int) (*UserModel, error) {
	return &UserModel{ID: id}, nil
}

func (m *MockUserStore) GetUserByEmail(context.Context, string) (*UserModel, error) {
	return &UserModel{}, nil
}

func (m *MockUserStore) CreateAndInvite(ctx context.Context, user *UserModel, token string, exp time.Duration) error {
	return nil
}

func (m *MockUserStore) createUserForInvitation(ctx context.Context, tx *sql.Tx, u *UserModel) error {
	return nil
}

func (m *MockUserStore) Activate(ctx context.Context, t string) error {
	return nil
}

func (m *MockUserStore) Delete(ctx context.Context, id int) error {
	return nil
}

func (m *MockUserStore) DeleteUserByID(ctx context.Context, id int) error {
	return nil
}

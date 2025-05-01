package cache

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/puremike/social-go/internal/store"
)

type UserStoreRdb struct {
	rdb *redis.Client
}

const timeExp = time.Minute * 2

func (r *UserStoreRdb) Get(ctx context.Context, id int) (*store.UserModel, error) {
	cacheKey := "user:" + strconv.Itoa(id)

	data, err := r.rdb.Get(ctx, cacheKey).Result()

	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var user store.UserModel

	if data != "" {
		if err := json.Unmarshal([]byte(data), &user); err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (r *UserStoreRdb) Set(ctx context.Context, user *store.UserModel) error {
	cacheKey := "user:" + strconv.Itoa(user.ID)

	json, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return r.rdb.SetEX(ctx, cacheKey, json, timeExp).Err()
}

func (r *UserStoreRdb) Delete(ctx context.Context, id int) {
	cacheKey := "user:" + strconv.Itoa(id)
	r.rdb.Del(ctx, cacheKey)
}

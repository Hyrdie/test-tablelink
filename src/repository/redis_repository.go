package repository

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"test-tablelink/src/entity"

	goRedis "github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *goRedis.Client
}

func NewRedisRepository(client *goRedis.Client) *RedisRepository {
	return &RedisRepository{client: client}
}

func (r *RedisRepository) SetUser(ctx context.Context, user *entity.User) error {
	key := "user:" + strconv.FormatInt(user.ID, 10)
	value, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, value, 24*time.Hour).Err()
}

func (r *RedisRepository) GetUser(ctx context.Context, id int64) (*entity.User, error) {
	key := "user:" + strconv.FormatInt(id, 10)
	value, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var user entity.User
	if err := json.Unmarshal(value, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *RedisRepository) DeleteUser(ctx context.Context, id int64) error {
	key := "user:" + strconv.FormatInt(id, 10)
	return r.client.Del(ctx, key).Err()
}

func (r *RedisRepository) SetToken(ctx context.Context, token string, userID int64) error {
	key := "token:" + token
	return r.client.Set(ctx, key, userID, 24*time.Hour).Err()
}

func (r *RedisRepository) GetUserIDByToken(ctx context.Context, token string) (int64, error) {
	key := "token:" + token
	userID, err := r.client.Get(ctx, key).Int64()
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (r *RedisRepository) DeleteToken(ctx context.Context, token string) error {
	key := "token:" + token
	return r.client.Del(ctx, key).Err()
}

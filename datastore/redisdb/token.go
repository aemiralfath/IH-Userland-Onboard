package redisdb

import (
	"context"
	"fmt"
	"time"

	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	"github.com/go-redis/redis/v8"
)

type TokenStore struct {
	redis *redis.Client
}

func NewTokenStore(redis *redis.Client) datastore.TokenStore {
	return &TokenStore{
		redis: redis,
	}
}

func (s *TokenStore) SetToken(ctx context.Context, tokenType, email, token string) error {
	key := fmt.Sprintf("token:%s:%s", tokenType, email)

	if err := s.redis.Set(ctx, key, token, time.Duration(time.Minute*5)); err.Err() != nil {
		return err.Err()
	}

	return nil
}

func (s *TokenStore) GetToken(ctx context.Context, tokenType, email string) (string, error) {
	key := fmt.Sprintf("token:%s:%s", tokenType, email)

	res := s.redis.Get(ctx, key)
	if res.Err() != nil {
		return "", res.Err()
	}

	token, err := res.Result()
	if err != nil {
		return "", err
	}

	return token, nil
}

package redisdb

import (
	"context"
	"fmt"
	"time"

	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	"github.com/go-redis/redis/v8"
)

type OTPStore struct {
	redis *redis.Client
}

func NewOTPStore(redis *redis.Client) datastore.OTPStore {
	return &OTPStore{
		redis: redis,
	}
}

func (s *OTPStore) GetOTP(ctx context.Context, email string, otp string) (string, error) {
	key := fmt.Sprintf("otp:%s", email)

	if err := s.redis.Set(ctx, key, otp, time.Duration(time.Hour*1)); err.Err() != nil {
		return "", err.Err()
	}

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

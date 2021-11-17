package redisdb

import (
	"context"
	"fmt"
	"os"
	"strconv"
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

func (s *OTPStore) SetOTP(ctx context.Context, otpType, otpCode, otpValue string) error {
	key := fmt.Sprintf("otp:%s:%s", otpType, otpCode)

	duration, err := strconv.Atoi(os.Getenv("OTP_DURATION"))
	if err != nil {
		return err
	}

	if err := s.redis.Set(ctx, key, otpValue, time.Duration(time.Minute*time.Duration(duration))); err.Err() != nil {
		return err.Err()
	}

	return nil
}

func (s *OTPStore) GetOTP(ctx context.Context, otpType, otpCode string) (string, error) {
	key := fmt.Sprintf("otp:%s:%s", otpType, otpCode)

	res := s.redis.Get(ctx, key)
	if res.Err() != nil {
		return "", res.Err()
	}

	otpValue, err := res.Result()
	if err != nil {
		return "", err
	}

	return otpValue, nil
}

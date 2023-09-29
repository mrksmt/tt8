package client

import (
	"context"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"github.com/mrksmt/tt8/service"
)

const redisEventKey = "redis_event_key"

// NewRedisRLClient returns new rate limited client with RedisDB limiter.
func NewRedisRLClient(
	rc *redis.Client,
	s service.Service,
) *Client {

	n, p := s.GetLimits()

	c := &Client{
		srv:     s,
		limiter: newRedisRateLimiter(rc, n, p),
	}

	return c
}

type RedisRateLimiter struct {
	limiter *redis_rate.Limiter
	rate    uint64
	period  time.Duration
}

var _ limiter = (*RedisRateLimiter)(nil)

func newRedisRateLimiter(
	rdb *redis.Client,
	rate uint64,
	period time.Duration,
) *RedisRateLimiter {
	rrl := &RedisRateLimiter{
		limiter: redis_rate.NewLimiter(rdb),
		rate:    rate,
		period:  period,
	}
	return rrl
}

func (rrl *RedisRateLimiter) AllowAtMost(
	ctx context.Context,
	n int,
) (
	allowed int,
	retryAfter time.Duration,
	err error,
) {

	// get rate limiting data from redis
	redisRateResult, err := rrl.limiter.AllowAtMost(
		ctx,
		redisEventKey,
		redis_rate.Limit{
			Rate:   int(rrl.rate),
			Burst:  int(rrl.rate),
			Period: rrl.period,
		},
		n,
	)

	if err != nil {
		return 0, 0, errors.Wrapf(ErrInternal, "redis ratelimit: %s", err)
	}

	return redisRateResult.Allowed, redisRateResult.RetryAfter, nil
}

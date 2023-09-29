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
// Implements GCRA (aka leaky bucket).
func NewRedisRLClient(
	rdb *redis.Client,
	srv service.Service,
) *Client {

	rate, period := srv.GetLimits()
	limiter := newRedisRateLimiter(rdb, rate, period)

	c := &Client{
		srv:     srv,
		limiter: limiter,
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
	limiter := redis_rate.NewLimiter(rdb)
	rrl := &RedisRateLimiter{
		limiter: limiter,
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

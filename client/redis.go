package client

import (
	"context"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

const redisEventKey = "redis_event_key"

type RedisRateLimiter struct {
	limiter *redis_rate.Limiter
}

var _ limiter = (*RedisRateLimiter)(nil)

func NewRedisRateLimiter(
	rc *redis.Client,
) *RedisRateLimiter {
	rrl := &RedisRateLimiter{limiter: redis_rate.NewLimiter(rc)}
	return rrl
}

func (rrl *RedisRateLimiter) AllowAtMost(
	ctx context.Context,
	rate int,
	period time.Duration,
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
			Rate:   int(n),
			Burst:  int(n),
			Period: period,
		},
		n,
	)

	if err != nil {
		return 0, 0, errors.Wrapf(ErrInternal, "redis ratelimit: %s", err)
	}

	return redisRateResult.Allowed, redisRateResult.RetryAfter, nil
}

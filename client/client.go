package client

import (
	"context"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/mrksmt/tt8/service"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

const redisEventKey = "redis_event_key"

// RedisRLClient rate limited client Redis ratelimit based.
type RedisRLClient struct {
	s       service.Service
	limiter *redis_rate.Limiter
}

var _ LimitedClient = (*RedisRLClient)(nil)

// NewRedisRLClient returns new RedisRLClient.
func NewRedisRLClient(
	rc *redis.Client,
	s service.Service,
) *RedisRLClient {
	limiter := redis_rate.NewLimiter(rc)

	c := &RedisRLClient{
		s:       s,
		limiter: limiter,
	}
	return c
}

// Process try ot process items.
// Returns number of processed items and retry duration.
func (c *RedisRLClient) Process(
	ctx context.Context,
	items ...service.Item,
) (
	processed int,
	retryAfter time.Duration,
	err error,
) {

	// check input length
	if len(items) == 0 {
		return 0, 0, nil
	}

	// get service limits (why not?)
	n, p := c.s.GetLimits()

	// get rate limiting data from redis
	redisRateResult, err := c.limiter.AllowAtMost(
		ctx,
		redisEventKey,
		redis_rate.Limit{
			Rate:   int(n),
			Burst:  int(n),
			Period: p,
		},
		len(items),
	)

	if err != nil {
		return 0, 0, errors.Wrap(err, "redis ratelimit")
	}

	// process allowed number of items
	err = c.s.Process(ctx, items[:redisRateResult.Allowed])
	if err != nil {
		return 0, 0, errors.Wrap(err, "process items")
	}

	return redisRateResult.Allowed, redisRateResult.RetryAfter, nil
}

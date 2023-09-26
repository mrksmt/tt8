package client

import (
	"context"
	"time"

	"github.com/mrksmt/tt8/service"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type limiter interface {
	AllowAtMost(
		ctx context.Context,
		rate int,
		period time.Duration,
		n int,
	) (
		allowed int,
		retryAfter time.Duration,
		err error,
	)
}

// ErrInternal any rate limiting error.
var ErrInternal = errors.New("internal")

// RedisRLClient rate limited client Redis ratelimit based.
type RedisRLClient struct {
	s       service.Service
	limiter limiter
}

var _ LimitedClient = (*RedisRLClient)(nil)

// NewRedisRLClient returns new RedisRLClient.
func NewRedisRLClient(
	rc *redis.Client,
	s service.Service,
) *RedisRLClient {

	c := &RedisRLClient{
		s:       s,
		limiter: NewRedisRateLimiter(rc),
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
	allowed, retryAfter, err := c.limiter.AllowAtMost(ctx, int(n), p, len(items))
	if err != nil {
		return 0, 0, errors.Wrap(err, "rate limit processing")
	}

	if allowed == 0 {
		return 0, retryAfter, nil
	}

	// add piece of paranoia
	if allowed > len(items) {
		allowed = len(items)
	}

	// process allowed number of items
	err = c.s.Process(ctx, items[:allowed])
	if err != nil {
		return 0, 0, errors.Wrap(err, "items processing")
	}

	return allowed, retryAfter, nil
}

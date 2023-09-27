package client

import (
	"context"
	"time"

	"github.com/mrksmt/tt8/service"
	"github.com/pkg/errors"
)

type limiter interface {
	AllowAtMost(
		ctx context.Context,
		n int,
	) (
		allowed int,
		retryAfter time.Duration,
		err error,
	)
}

// ErrInternal any rate limiting error.
var ErrInternal = errors.New("internal")

// Client rate limited client Redis ratelimit based.
type Client struct {
	srv     service.Service
	limiter limiter
}

var _ LimitedClient = (*Client)(nil)

// Process implements LimitedClient interface
func (c *Client) Process(
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

	// get rate limiting data from limiter
	allowed, retryAfter, err := c.limiter.AllowAtMost(ctx, len(items))
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
	err = c.srv.Process(ctx, items[:allowed])
	if err != nil {
		return 0, 0, errors.Wrap(err, "items processing")
	}

	return allowed, retryAfter, nil
}

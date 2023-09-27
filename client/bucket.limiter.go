package client

import (
	"context"
	"sync"
	"time"

	"github.com/mrksmt/tt8/service"
)

// NewTokenBucketClient returns new rate limited client with token bucket limiter.
func NewTokenBucketClient(
	srv service.Service,
) *Client {

	rate, period := srv.GetLimits()

	c := &Client{
		srv:     srv,
		limiter: newBucketRateLimiter(int(rate), period),
	}

	return c
}

type BucketRateLimiter struct {
	rate      int
	period    time.Duration
	mx        sync.RWMutex
	tokens    int
	updatedAt time.Time
}

var _ limiter = (*BucketRateLimiter)(nil)

func newBucketRateLimiter(
	rate int,
	period time.Duration,
) *BucketRateLimiter {

	brl := &BucketRateLimiter{
		rate:   rate,
		period: period,
	}

	return brl
}

func (brl *BucketRateLimiter) AllowAtMost(
	ctx context.Context,
	n int,
) (
	allowed int,
	retryAfter time.Duration,
	err error,
) {

	brl.mx.RLock()
	defer brl.mx.RUnlock()

	retryAfter = -1

	if time.Since(brl.updatedAt) >= brl.period {
		brl.tokens = brl.rate
		brl.updatedAt = time.Now()
	}

	allowed = brl.tokens
	brl.tokens -= n

	if brl.tokens < 0 {
		brl.tokens = 0
	}

	if brl.tokens == 0 {
		retryAfter = brl.period - time.Since(brl.updatedAt)
	}

	if brl.tokens > 0 {
		allowed = n
	}

	return
}

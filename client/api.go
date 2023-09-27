package client

import (
	"context"
	"time"

	"github.com/mrksmt/tt8/service"
)

// LimitedClient defines common interface for all kinds of rate limited clients.
// All clients should be implements it.
// processed - number of success processed items.
// retryAfter - until the next request will be permitted.
type LimitedClient interface {
	Process(
		ctx context.Context,
		items ...service.Item,
	) (
		processed int,
		retryAfter time.Duration,
		err error,
	)
}

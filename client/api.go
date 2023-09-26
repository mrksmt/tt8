package client

import (
	"context"
	"time"

	"github.com/mrksmt/tt8/service"
)

// LimitedClient ic common client interface.
// All clients should be implements it.
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

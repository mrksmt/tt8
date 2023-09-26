package client

import (
	"context"
	"testing"
	"time"

	"github.com/mrksmt/tt8/client/mocks"
	"github.com/mrksmt/tt8/service"
	"github.com/stretchr/testify/require"
)

func Test_RateLimit(t *testing.T) {

	mockedService := &mocks.MockedService{}
	mockedService.On("GetLimits").Return(uint64(10), time.Minute)
	mockedService.On("Process").Return(nil)

	mockedRedisLimiter := &mocks.MockedRedisLimiter{}
	mockedRedisLimiter.On("AllowAtMost").Return(2, time.Duration(-1), error(nil))

	client := RedisRLClient{
		s:       mockedService,
		limiter: mockedRedisLimiter,
	}

	processed, retryAfter, err := client.Process(
		context.TODO(),
		make([]service.Item, 3)...,
	)

	require.NoError(t, err)
	require.LessOrEqual(t, retryAfter, time.Duration(0))
	require.EqualValues(t, 2, processed)
}

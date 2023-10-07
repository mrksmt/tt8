package client

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"

	"github.com/mrksmt/tt8/client/mocks"
	"github.com/mrksmt/tt8/service"
)

func Test_RateLimit(t *testing.T) {

	mockedService := &mocks.MockedService{}
	mockedService.On("GetLimits").Return(uint64(10), time.Minute)
	mockedService.On("Process").Return(nil)

	mockedRedisLimiter := &mocks.MockedRedisLimiter{}
	mockedRedisLimiter.On("AllowAtMost").Return(2, time.Duration(-1), error(nil))

	client := Client{
		srv:     mockedService,
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

func Test_redisLimiter(t *testing.T) {

	ctx := context.TODO()

	redisOptions := &redis.Options{
		Addr: "localhost:46379",
		DB:   1,
	}

	rdb := redis.NewClient(redisOptions)
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Skip(err)
	}
	if err := rdb.FlushAll(ctx).Err(); err != nil {
		t.Skip(err)
	}

	limiter := newRedisRateLimiter(rdb, 8, time.Second)

	tests := []struct {
		name  string
		n     int
		allow int
		sleep time.Duration
	}{
		{
			name:  "1",
			n:     2,
			allow: 2,
			sleep: time.Millisecond * 100,
		},
		{
			name:  "2",
			n:     2,
			allow: 2,
			sleep: time.Millisecond * 100,
		},
		{
			name:  "3",
			n:     8,
			allow: 5,
			sleep: time.Millisecond * 100,
		},
		{
			name:  "4",
			n:     1,
			allow: 0,
			sleep: time.Millisecond * 1000,
		},
		{
			name:  "5",
			n:     8,
			allow: 8,
			sleep: time.Millisecond * 100,
		},
		{
			name:  "6",
			n:     1,
			allow: 0,
		},
	}

	for _, test := range tests {

		testFunc := func(t *testing.T) {
			allowed, _, err := limiter.AllowAtMost(ctx, test.n)
			require.NoError(t, err)
			require.Equal(t, test.allow, allowed)
		}

		t.Run(test.name, testFunc)
		<-time.After(test.sleep)
	}
}

func Test_redisLimiter2(t *testing.T) {

	ctx := context.TODO()

	redisOptions := &redis.Options{
		Addr: "localhost:46379",
		DB:   1,
	}

	rdb := redis.NewClient(redisOptions)
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Skip(err)
	}
	if err := rdb.FlushDB(ctx).Err(); err != nil {
		t.Skip(err)
	}

	rdb.PubSubShardChannels(ctx, "")
}

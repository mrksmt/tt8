package tt8_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/mrksmt/tt8/client"
	"github.com/mrksmt/tt8/client/mocks"
	"github.com/mrksmt/tt8/service"
)

func getMockedService() service.Service {
	mockedService := &mocks.MockedService{}
	mockedService.On("GetLimits").Return(uint64(8), time.Second)
	mockedService.On("Process").Return(nil)
	return mockedService
}

// Test_ExampleRedisClient run ExampleRedisClient func
func Test_ExampleRedisClient(t *testing.T) {
	ExampleBucketClient()
}

func ExampleRedisClient() {

	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*2000)
	defer cancel()

	redisOptions := &redis.Options{
		Addr: "localhost:46379",
		DB:   1,
	}

	rdb := redis.NewClient(redisOptions)
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal(err)
	}
	if err := rdb.FlushAll(ctx).Err(); err != nil {
		log.Fatal(err)
	}

	limitedClient := client.NewRedisRLClient(
		rdb,
		getMockedService(),
	)

	start := time.Now()
	ticker := time.NewTicker(time.Millisecond * 100)

	for {

		batch := make(service.Batch, 3)

		processed, retryAfter, err := limitedClient.Process(ctx, batch...)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf(
			"since start: %s\nitems processed: %d\nretry after: %s\n\n",
			time.Since(start).Round(time.Millisecond),
			processed,
			retryAfter.Round(time.Millisecond),
		)

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}

}

// Test_ExampleBucketClient run ExampleBucketClient func
func Test_ExampleBucketClient(t *testing.T) {
	ExampleBucketClient()
}

// ExampleBucketClient example run of bucket rate limited client
func ExampleBucketClient() {

	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*2000)
	defer cancel()

	limitedClient := client.NewTokenBucketClient(
		getMockedService(),
	)

	start := time.Now()
	ticker := time.NewTicker(time.Millisecond * 100)

	for {

		batch := make(service.Batch, 3)

		processed, retryAfter, err := limitedClient.Process(ctx, batch...)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf(
			"since start: %s\nitems processed: %d\nretry after: %s\n\n",
			time.Since(start).Round(time.Millisecond),
			processed,
			retryAfter.Round(time.Millisecond),
		)

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

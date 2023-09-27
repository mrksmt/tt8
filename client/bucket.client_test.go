package client

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_bucketLimiter(t *testing.T) {

	ctx := context.TODO()

	limiter := newBucketRateLimiter(8, time.Second)

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
			allow: 4,
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

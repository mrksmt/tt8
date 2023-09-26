package mocks

import (
	"context"
	"time"

	"github.com/mrksmt/tt8/service"
	"github.com/stretchr/testify/mock"
)

type MockedRedisLimiter struct{ mock.Mock }

func (s *MockedRedisLimiter) AllowAtMost(
	ctx context.Context,
	rate int,
	period time.Duration,
	n int,
) (
	allowed int,
	retryAfter time.Duration,
	err error,
) {
	args := s.Called()

	err, ok := args.Get(2).(error)
	_ = ok

	return args.Get(0).(int), args.Get(1).(time.Duration), err
}

type MockedService struct{ mock.Mock }

func (s *MockedService) GetLimits() (n uint64, p time.Duration) {
	args := s.Called()
	return args.Get(0).(uint64), args.Get(1).(time.Duration)
}

func (s *MockedService) Process(ctx context.Context, batch service.Batch) error {
	args := s.Called()

	err, ok := args.Get(0).(error)
	_ = ok

	return err
}

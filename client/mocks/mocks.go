package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/mrksmt/tt8/service"
)

type MockedRedisLimiter struct{ mock.Mock }

func (s *MockedRedisLimiter) AllowAtMost(
	ctx context.Context,
	n int,
) (
	allowed int,
	retryAfter time.Duration,
	err error,
) {
	args := s.Called()

	err, ok := args.Get(2).(error)
	_ = ok

	return args.Get(0).(int), args.Get(1).(time.Duration), err //nolint
}

type MockedService struct{ mock.Mock }

func (s *MockedService) GetLimits() (n uint64, p time.Duration) {
	args := s.Called()
	return args.Get(0).(uint64), args.Get(1).(time.Duration) //nolint
}

func (s *MockedService) Process(ctx context.Context, batch service.Batch) error {
	args := s.Called()

	err, ok := args.Get(0).(error)
	_ = ok

	return err
}

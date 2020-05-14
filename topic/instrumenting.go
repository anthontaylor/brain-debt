package topic

import (
	"context"
	"time"

	brain "github.com/anthontaylor/brain-debt"

	"github.com/go-kit/kit/metrics"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

func NewInstrumentingService(counter metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		Service:        s,
	}
}

func (s *instrumentingService) Add(ctx context.Context, id *brain.UserID, name string) (*brain.TopicID, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "add").Add(1)
		s.requestLatency.With("method", "add").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.Add(ctx, id, name)
}

func (s *instrumentingService) Get(ctx context.Context, id *brain.UserID) (user []brain.Topic, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "find").Add(1)
		s.requestLatency.With("method", "find").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.Get(ctx, id)
}

func (s *instrumentingService) Update(ctx context.Context, id *brain.UserID, topic *brain.Topic) (_ *brain.Topic, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "update").Add(1)
		s.requestLatency.With("method", "update").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.Update(ctx, id, topic)
}

func (s *instrumentingService) Delete(ctx context.Context, id *brain.UserID, topicID *brain.TopicID) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "delete").Add(1)
		s.requestLatency.With("method", "delete").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.Delete(ctx, id, topicID)
}

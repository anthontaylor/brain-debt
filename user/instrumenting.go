package user

import (
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

func (s *instrumentingService) Add(user brain.User) (*brain.UserID, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "add").Add(1)
		s.requestLatency.With("method", "add").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.Add(user)
}

func (s *instrumentingService) Find(id *brain.UserID) (user *brain.User, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "find").Add(1)
		s.requestLatency.With("method", "find").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.Find(id)
}

func (s *instrumentingService) Update(id *brain.UserID, user brain.User) (_ *brain.User, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "update").Add(1)
		s.requestLatency.With("method", "update").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.Update(id, user)
}

func (s *instrumentingService) Delete(id *brain.UserID) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "delete").Add(1)
		s.requestLatency.With("method", "delete").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.Delete(id)
}

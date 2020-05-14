package topic

import (
	"context"
	"time"

	brain "github.com/anthontaylor/brain-debt"
	"github.com/go-kit/kit/log"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger: logger, Service: s}
}

func (s *loggingService) Add(ctx context.Context, id *brain.UserID, name string) (_ *brain.TopicID, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "add",
			"topic name", name,
			"took", time.Since(begin),
			"id", id,
			"err", err,
		)
	}(time.Now())
	return s.Service.Add(ctx, id, name)
}

func (s *loggingService) Get(ctx context.Context, id *brain.UserID) (_ []brain.Topic, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "find",
			"took", time.Since(begin),
			"id", id,
			"err", err,
		)
	}(time.Now())
	return s.Service.Get(ctx, id)
}

func (s *loggingService) Update(ctx context.Context, id *brain.UserID, topic *brain.Topic) (_ *brain.Topic, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "update",
			"took", time.Since(begin),
			"user id", id,
			"topic id", topic.ID,
			"topic name", topic.Name,
			"err", err,
		)
	}(time.Now())
	return s.Service.Update(ctx, id, topic)
}

func (s *loggingService) Delete(ctx context.Context, id *brain.UserID, topicID *brain.TopicID) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "delete",
			"took", time.Since(begin),
			"user id", id,
			"topic id", topicID,
			"err", err,
		)
	}(time.Now())
	return s.Service.Delete(ctx, id, topicID)
}

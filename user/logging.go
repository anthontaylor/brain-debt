package user

import (
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

func (s *loggingService) Add(f, l string) (id brain.UserID, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "add",
			"firstName", f,
			"lastName", l,
			"took", time.Since(begin),
			"id", id,
			"err", err,
		)
	}(time.Now())
	return s.Service.Add(f, l)
}

func (s *loggingService) Find(id brain.UserID) (user *brain.User, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "find",
			"took", time.Since(begin),
			"id", id,
			"err", err,
		)
	}(time.Now())
	return s.Service.Find(id)
}

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

func (s *loggingService) Add(user brain.User) (id *brain.UserID, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "add",
			"firstName", user.FirstName,
			"lastName", user.LastName,
			"took", time.Since(begin),
			"id", id,
			"err", err,
		)
	}(time.Now())
	return s.Service.Add(user)
}

func (s *loggingService) Find(id *brain.UserID) (user *brain.User, err error) {
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

func (s *loggingService) Update(id *brain.UserID, user brain.User) (_ *brain.User, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "update",
			"took", time.Since(begin),
			"id", id,
			"firstName", user.FirstName,
			"lastName", user.LastName,
			"err", err,
		)
	}(time.Now())
	return s.Service.Update(id, user)
}

func (s *loggingService) Delete(id *brain.UserID) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "delete",
			"took", time.Since(begin),
			"id", id,
			"err", err,
		)
	}(time.Now())
	return s.Service.Delete(id)
}

package topic

import (
	"context"
	"errors"

	brain "github.com/anthontaylor/brain-debt"
)

var ErrInvalidArgument = errors.New("invalid argument")

//go:generate mockgen -destination=mocks/service.go -package=topic_mock github.com/anthontaylor/brain-debt/topic Service

type Service interface {
	Add(ctx context.Context, id *brain.UserID, tName string) (*brain.TopicID, error)
	Get(ctx context.Context, id *brain.UserID) ([]brain.Topic, error)
	Update(ctx context.Context, id *brain.UserID, topic *brain.Topic) (*brain.Topic, error)
	Delete(ctx context.Context, id *brain.UserID, topicID *brain.TopicID) error
}

type service struct {
	topic brain.TopicRepository
}

func NewService(topic brain.TopicRepository) Service {
	return &service{topic: topic}
}

func (s *service) Add(ctx context.Context, id *brain.UserID, name string) (*brain.TopicID, error) {
	if id == nil || *id == "" || name == "" {
		return nil, ErrInvalidArgument
	}
	return s.topic.Add(ctx, id, name)
}

func (s *service) Get(ctx context.Context, id *brain.UserID) ([]brain.Topic, error) {
	if id == nil || *id == "" {
		return nil, ErrInvalidArgument
	}
	return s.topic.Get(ctx, id)
}

func (s *service) Update(ctx context.Context, id *brain.UserID, topic *brain.Topic) (*brain.Topic, error) {
	if id == nil || topic == nil {
		return nil, ErrInvalidArgument
	}
	return s.topic.Update(ctx, id, topic)
}

func (s *service) Delete(ctx context.Context, id *brain.UserID, topicID *brain.TopicID) error {
	if id == nil || topicID == nil {
		return ErrInvalidArgument
	}
	return s.topic.Delete(ctx, id, topicID)
}

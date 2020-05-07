package topic

import (
	"errors"

	brain "github.com/anthontaylor/brain-debt"
)

var ErrInvalidArgument = errors.New("invalid argument")

type Service interface {
	Add(id *brain.UserID, tName string) (*brain.TopicID, error)
	Get(id *brain.UserID) ([]brain.Topic, error)
	Update(id *brain.UserID, topic *brain.Topic) (*brain.Topic, error)
	Delete(id *brain.UserID, topicID *brain.TopicID) error
}

type service struct {
	topic brain.TopicRepository
}

func NewService(topic brain.TopicRepository) Service {
	return &service{topic: topic}
}

func (s *service) Add(id *brain.UserID, name string) (*brain.TopicID, error) {
	if id == nil || name == "" {
		return nil, ErrInvalidArgument
	}
	return s.topic.Add(id, name)
}

func (s *service) Get(id *brain.UserID) ([]brain.Topic, error) {
	if id == nil {
		return nil, ErrInvalidArgument
	}
	return s.topic.Get(id)
}

func (s *service) Update(id *brain.UserID, topic *brain.Topic) (*brain.Topic, error) {
	if id == nil || topic == nil {
		return nil, ErrInvalidArgument
	}
	return s.topic.Update(id, topic)
}

func (s *service) Delete(id *brain.UserID, topicID *brain.TopicID) error {
	if id == nil || topicID == nil {
		return ErrInvalidArgument
	}
	return s.topic.Delete(id, topicID)
}

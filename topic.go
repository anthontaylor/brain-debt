package brain

import (
	"errors"
)

type Topic struct {
	ID   TopicID `json:"id"`
	Name string  `json:"name"`
}

type TopicID string

type TopicRepository interface {
	Add(id *UserID, tName string) (*TopicID, error)
	Get(id *UserID) ([]Topic, error)
	Update(id *UserID, topic *Topic) (*Topic, error)
	Delete(id *UserID, topicID *TopicID) error
}

var ErrTopicNotFound = errors.New("topic not found")

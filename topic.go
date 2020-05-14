package brain

import (
	"context"
	"errors"
)

type Topic struct {
	ID   TopicID `json:"id"`
	Name string  `json:"name"`
}

type TopicID string

type TopicRepository interface {
	Add(ctx context.Context, id *UserID, tName string) (*TopicID, error)
	Get(ctx context.Context, id *UserID) ([]Topic, error)
	Update(ctx context.Context, id *UserID, topic *Topic) (*Topic, error)
	Delete(ctx context.Context, id *UserID, topicID *TopicID) error
}

var ErrTopicNotFound = errors.New("topic not found")

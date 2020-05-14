package inmem

import (
	"context"
	"sync"

	brain "github.com/anthontaylor/brain-debt"
	"github.com/google/uuid"
)

type topicRepository struct {
	mu     sync.RWMutex
	topics map[brain.UserID][]brain.Topic
}

func NewTopicRepository() brain.TopicRepository {
	return &topicRepository{
		topics: make(map[brain.UserID][]brain.Topic),
	}
}

func (r *topicRepository) Add(ctx context.Context, id *brain.UserID, tName string) (*brain.TopicID, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	topicID := brain.TopicID(uuid.New().String())
	if topics, ok := r.topics[*id]; ok {
		r.topics[*id] = append(topics, brain.Topic{ID: topicID, Name: tName})
		return &topicID, nil

	}
	r.topics[*id] = []brain.Topic{
		brain.Topic{ID: topicID, Name: tName},
	}
	return &topicID, nil
}

func (r *topicRepository) Get(ctx context.Context, id *brain.UserID) ([]brain.Topic, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	topics := r.topics[*id]
	return topics, nil
}

func (r *topicRepository) Update(ctx context.Context, id *brain.UserID, topic *brain.Topic) (*brain.Topic, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if topics, ok := r.topics[*id]; ok {
		for i, t := range topics {
			if t.ID == topic.ID {
				t.Name = topic.Name
				topics[i] = t
			}
		}
		return topic, nil
	}
	return nil, brain.ErrTopicNotFound
}

func (r *topicRepository) Delete(ctx context.Context, id *brain.UserID, topicID *brain.TopicID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if topics, ok := r.topics[*id]; ok {
		for i, topic := range topics {
			if *topicID == topic.ID {
				topics = removeTopic(topics, i)
				r.topics[*id] = topics
				return nil
			}
		}
		return brain.ErrTopicNotFound
	}
	return brain.ErrUnknownUser
}

func removeTopic(s []brain.Topic, i int) []brain.Topic {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

package cassandra

import (
	"context"
	"fmt"

	brain "github.com/anthontaylor/brain-debt"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

type topicRepository struct {
	conn *gocql.Session
}

func NewTopicRepository(conn *gocql.Session) brain.TopicRepository {
	return &topicRepository{conn: conn}
}

func (r *topicRepository) Add(ctx context.Context, id *brain.UserID, tName string) (*brain.TopicID, error) {
	userGuid, err := gocql.ParseUUID(string(*id))
	if err != nil {
		return nil, fmt.Errorf("user_guid=%v %w", *id, brain.ErrInvalidUserID)
	}
	topicID := uuid.New().String()
	topicGuid, err := gocql.ParseUUID(topicID)
	if err != nil {
		return nil, fmt.Errorf("user_guid=%v %w", *id, brain.ErrInvalidTopicID)
	}
	topicMap := map[gocql.UUID]string{topicGuid: tName}
	query := r.conn.Query(`UPDATE topic SET topics = topics + ? WHERE user_guid = ?`).
		Consistency(gocql.LocalQuorum).
		Bind(topicMap, userGuid)
	if err := query.Exec(); err != nil {
		return nil, fmt.Errorf("user_guid=%v %w err=%v", *id, brain.ErrUnknownUser, err)
	}
	tID := brain.TopicID(topicID)
	return &tID, nil
}

func (r *topicRepository) Get(ctx context.Context, id *brain.UserID) ([]brain.Topic, error) {
	userGuid, err := gocql.ParseUUID(string(*id))
	if err != nil {
		return nil, fmt.Errorf("user_guid=%v %w", *id, brain.ErrInvalidUserID)
	}
	query := r.conn.Query(`SELECT topics FROM topic WHERE user_guid = ?`, userGuid).
		Consistency(gocql.LocalQuorum).
		Bind(userGuid)

	topicsMap := make(map[gocql.UUID]string)
	if err := query.Scan(&topicsMap); err != nil {
		return nil, fmt.Errorf("user_guid=%v %w err=%v", *id, brain.ErrUnknownUser, err)
	}
	var topics []brain.Topic
	for k, v := range topicsMap {
		topic := brain.Topic{ID: brain.TopicID(k.String()), Name: v}
		topics = append(topics, topic)
	}
	return topics, nil
}

func (r *topicRepository) Update(ctx context.Context, id *brain.UserID, topic *brain.Topic) (*brain.Topic, error) {
	userGuid, err := gocql.ParseUUID(string(*id))
	if err != nil {
		return nil, fmt.Errorf("user_guid=%v %w", *id, brain.ErrInvalidUserID)
	}
	topicGuid, err := gocql.ParseUUID(string(topic.ID))
	if err != nil {
		return nil, fmt.Errorf("user_guid=%v %w", *id, brain.ErrInvalidTopicID)
	}
	query := r.conn.Query(`UPDATE topic SET topics[?] = ? WHERE user_guid = ?`).
		Consistency(gocql.LocalQuorum).
		Bind(topicGuid, topic.Name, userGuid)
	if err := query.Exec(); err != nil {
		return nil, fmt.Errorf("user_guid=%v %w err=%v", *id, brain.ErrUnknownUser, err)
	}
	return topic, nil
}

func (r *topicRepository) Delete(ctx context.Context, id *brain.UserID, topicID *brain.TopicID) error {
	userGuid, err := gocql.ParseUUID(string(*id))
	if err != nil {
		return fmt.Errorf("user_guid=%v %w", *id, brain.ErrInvalidUserID)
	}
	topicGuid, err := gocql.ParseUUID(string(*topicID))
	if err != nil {
		return fmt.Errorf("user_guid=%v %w", *id, brain.ErrInvalidTopicID)
	}
	query := r.conn.Query("DELETE topics[?] FROM topic WHERE user_guid = ?").
		Consistency(gocql.LocalQuorum).
		Bind(topicGuid, userGuid)
	if err := query.Exec(); err != nil {
		return fmt.Errorf("user_guid=%v %w err=%v", *id, brain.ErrUnknownUser, err)
	}
	return nil
}

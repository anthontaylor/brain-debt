package topic

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	brain "github.com/anthontaylor/brain-debt"
)

type addTopicRequest struct {
	ID   *brain.UserID
	Name string
}

type addTopicResponse struct {
	ID  *brain.TopicID `json:"topic_id,omitempty"`
	Err error          `json:"error,omitempty"`
}

func makeAddTopicEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addTopicRequest)
		id, err := s.Add(req.ID, req.Name)
		if err != nil {
			return addTopicResponse{Err: err}, err
		}
		return addTopicResponse{ID: id, Err: err}, nil
	}
}

type getTopicsRequest struct {
	ID *brain.UserID
}

type getTopicsResponse struct {
	Topics []brain.Topic `json:"topics,omitempty"`
	Err    error         `json:"error,omitempty"`
}

func makeGetTopicsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getTopicsRequest)
		topics, err := s.Get(req.ID)
		if err != nil {
			return getTopicsResponse{Err: err}, err
		}
		return getTopicsResponse{Topics: topics, Err: err}, nil
	}
}

type updateTopicRequest struct {
	UserID *brain.UserID
	Topic  *brain.Topic
}

type updateTopicResponse struct {
	Topic *brain.Topic `json:"topic,omitempty"`
	Err   error        `json:"error,omitempty"`
}

func makeUpdateTopicEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateTopicRequest)
		topic, err := s.Update(req.UserID, req.Topic)
		if err != nil {
			return updateTopicResponse{Err: err}, err
		}
		return updateTopicResponse{Topic: topic, Err: err}, nil
	}
}

type deleteTopicRequest struct {
	UserID  *brain.UserID
	TopicID *brain.TopicID
}

type deleteTopicResponse struct {
	Err error `json:"error,omitempty"`
}

func makeDeleteTopicEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteTopicRequest)
		err := s.Delete(req.UserID, req.TopicID)
		if err != nil {
			return deleteTopicResponse{Err: err}, err
		}
		return deleteTopicResponse{Err: err}, nil
	}
}

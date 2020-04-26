package user

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	brain "github.com/anthontaylor/brain-debt"
)

type addUserRequest struct {
	FirstName string
	LastName  string
}

type addUserResponse struct {
	ID  brain.UserID `json:"user_id,omitempty"`
	Err error        `json:"error,omitempty"`
}

func makeAddUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addUserRequest)
		id, err := s.Add(req.FirstName, req.LastName)
		return addUserResponse{ID: id, Err: err}, nil
	}
}

type getUserRequest struct {
	ID brain.UserID
}

type getUserResponse struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Err       error  `json:"error,omitempty"`
}

func makeGetUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getUserRequest)
		user, err := s.Find(req.ID)
		if err != nil {
			return getUserResponse{Err: err}, nil
		}
		return getUserResponse{FirstName: user.FirstName, LastName: user.LastName, Err: err}, nil
	}
}

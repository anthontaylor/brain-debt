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
	ID  *brain.UserID `json:"user_id,omitempty"`
	Err error         `json:"error,omitempty"`
}

func makeAddUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addUserRequest)
		id, err := s.Add(brain.User{req.FirstName, req.LastName})
		if err != nil {
			return addUserResponse{Err: err}, err
		}
		return addUserResponse{ID: id, Err: err}, nil
	}
}

type getUserRequest struct {
	ID *brain.UserID
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
			return getUserResponse{Err: err}, err
		}
		return getUserResponse{FirstName: user.FirstName, LastName: user.LastName, Err: err}, nil
	}
}

type updateUserRequest struct {
	ID        *brain.UserID `json:"id,omitempty"`
	FirstName string        `json:"first_name,omitempty"`
	LastName  string        `json:"last_name,omitempty"`
}

type updateUserResponse struct {
	ID        *brain.UserID `json:"id",omitempty`
	FirstName string        `json:"first_name,omitempty"`
	LastName  string        `json:"last_name,omitempty"`
	Err       error         `json:"error,omitempty"`
}

func makeUpdateUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateUserRequest)
		user, err := s.Update(req.ID, brain.User{FirstName: req.FirstName, LastName: req.LastName})
		if err != nil {
			return updateUserResponse{Err: err}, err
		}
		return updateUserResponse{ID: req.ID, FirstName: user.FirstName, LastName: user.LastName, Err: err}, nil
	}
}

type deleteUserRequest struct {
	ID *brain.UserID `json:"id,omitempty"`
}

type deleteUserResponse struct {
	Err error `json:"error,omitempty"`
}

func makeDeleteUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteUserRequest)
		err := s.Delete(req.ID)
		if err != nil {
			return deleteUserResponse{Err: err}, err
		}
		return deleteUserResponse{}, nil

	}
}

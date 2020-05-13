package user

import (
	"context"
	"errors"

	brain "github.com/anthontaylor/brain-debt"
)

var ErrInvalidArgument = errors.New("invalid argument")

type Service interface {
	Add(ctx context.Context, user brain.User) (*brain.UserID, error)
	Find(ctx context.Context, id *brain.UserID) (*brain.User, error)
	Update(ctx context.Context, id *brain.UserID, user brain.User) (*brain.User, error)
	Delete(ctx context.Context, id *brain.UserID) error
}

type service struct {
	users brain.UserRepository
}

func NewService(users brain.UserRepository) Service {
	return &service{users: users}
}

func (s *service) Add(ctx context.Context, user brain.User) (*brain.UserID, error) {
	if user.FirstName == "" || user.LastName == "" {
		return nil, ErrInvalidArgument
	}
	return s.users.Add(ctx, &user)
}

func (s *service) Find(ctx context.Context, id *brain.UserID) (*brain.User, error) {
	if id == nil {
		return nil, ErrInvalidArgument
	}
	return s.users.Find(ctx, id)
}

func (s *service) Update(ctx context.Context, id *brain.UserID, user brain.User) (*brain.User, error) {
	if id == nil || user.FirstName == "" || user.LastName == "" {
		return nil, ErrInvalidArgument
	}
	return s.users.Update(ctx, id, user)
}

func (s *service) Delete(ctx context.Context, id *brain.UserID) error {
	if id == nil {
		return ErrInvalidArgument
	}
	return s.users.Delete(ctx, id)
}
